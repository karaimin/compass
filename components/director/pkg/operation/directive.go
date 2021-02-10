/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/pkg/header"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/internal/model"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

const ModeParam = "mode"

// ResourceFetcherFunc defines a function which fetches the webhooks for a specific resource ID
type WebhookFetcherFunc func(ctx context.Context, resourceID string) ([]*model.Webhook, error)

type directive struct {
	transact            persistence.Transactioner
	webhookFetcherFunc  WebhookFetcherFunc
	tenantLoaderFunc    TenantLoaderFunc
	resourceFetcherFunc ResourceFetcherFunc
	scheduler           Scheduler
}

// NewDirective creates a new handler struct responsible for the Async directive business logic
func NewDirective(transact persistence.Transactioner, webhookFetcherFunc WebhookFetcherFunc, resourceFetcherFunc ResourceFetcherFunc, tenantLoaderFunc TenantLoaderFunc, scheduler Scheduler) *directive {
	return &directive{
		transact:            transact,
		webhookFetcherFunc:  webhookFetcherFunc,
		tenantLoaderFunc:    tenantLoaderFunc,
		resourceFetcherFunc: resourceFetcherFunc,
		scheduler:           scheduler,
	}
}

// HandleOperation enriches the request with an Operation information when the requesting mutation is annotated with the Async directive
func (d *directive) HandleOperation(ctx context.Context, _ interface{}, next gqlgen.Resolver, operationType graphql.OperationType, webhookType graphql.WebhookType, idField, parentIdField *string) (res interface{}, err error) {
	resCtx := gqlgen.GetResolverContext(ctx)
	var mode graphql.OperationMode
	if _, found := resCtx.Args[ModeParam]; !found {
		mode = graphql.OperationModeSync
	} else {
		modePointer, ok := resCtx.Args[ModeParam].(*graphql.OperationMode)
		if !ok {
			return nil, apperrors.NewInternalError(fmt.Sprintf("could not get %s parameter", ModeParam))
		}
		mode = *modePointer
	}

	ctx = SaveModeToContext(ctx, mode)

	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while opening database transaction: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to initialize database operation")
	}
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)

	if err := d.concurrencyCheck(ctx, operationType, resCtx, idField, parentIdField); err != nil {
		return nil, err
	}

	if mode == graphql.OperationModeSync {
		resp, err := next(ctx)
		if err != nil {
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.C(ctx).WithError(err).Errorf("An error occurred while closing database transaction: %s", err.Error())
			return nil, apperrors.NewInternalError("Unable to finalize database operation")
		}

		return resp, nil
	}

	operation := &Operation{
		OperationType:     OperationType(operationType),
		OperationCategory: resCtx.Field.Name,
		CorrelationID:     log.C(ctx).Data[log.FieldRequestID].(string),
	}
	ctx = SaveToContext(ctx, &[]*Operation{operation})

	resp, err := next(ctx)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while processing operation: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to process operation")
	}

	entity, ok := resp.(graphql.Entity)
	if !ok {
		log.C(ctx).WithError(err).Error("An error occurred while casting the response entity")
		return nil, apperrors.NewInternalError("Failed to process operation")
	}

	operation.ResourceID = entity.GetID()
	operation.ResourceType = entity.GetType()

	webhookIDs, err := d.prepareWebhookIDs(ctx, err, operation, webhookType)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while retrieving webhooks: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to retrieve webhooks")
	}

	operation.WebhookIDs = webhookIDs

	requestData, err := d.prepareRequestData(ctx, err, resp)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while preparing request data: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to prepare webhook request data")
	}

	operation.RequestData = requestData

	operationID, err := d.scheduler.Schedule(ctx, operation)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while scheduling operation: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to schedule operation")
	}

	operation.OperationID = operationID

	err = tx.Commit()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error occurred while closing database transaction: %s", err.Error())
		return nil, apperrors.NewInternalError("Unable to finalize database operation")
	}

	return resp, nil
}

func (d *directive) concurrencyCheck(ctx context.Context, op graphql.OperationType, resCtx *gqlgen.ResolverContext, idField, parentIdField *string) error {
	if op == graphql.OperationTypeCreate && parentIdField == nil {
		return nil
	}

	if idField == nil && parentIdField == nil {
		return apperrors.NewInternalError("idField or parentIdField from context should not be empty")
	}

	var resourceID string
	var ok bool
	if parentIdField != nil {
		resourceID, ok = resCtx.Args[*parentIdField].(string)
		if !ok {
			return apperrors.NewInternalError(fmt.Sprintf("could not get parentIdField: %q from request context", *parentIdField))
		}
	} else {
		resourceID, ok = resCtx.Args[*idField].(string)
		if !ok {
			return apperrors.NewInternalError(fmt.Sprintf("could not get idField: %q from request context", *idField))
		}
	}

	tenant, err := d.tenantLoaderFunc(ctx)
	if err != nil {
		return apperrors.NewTenantRequiredError()
	}

	app, err := d.resourceFetcherFunc(ctx, tenant, resourceID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return err
		}

		return apperrors.NewInternalError("failed to fetch resource with id %s", resourceID)
	}

	if app.GetDeletedAt().IsZero() && app.GetUpdatedAt().IsZero() && !app.GetReady() && (app.GetError() == nil || *app.GetError() == "") { // CREATING
		return apperrors.NewConcurrentOperationInProgressError("create operation is in progress")
	}
	if !app.GetDeletedAt().IsZero() && (app.GetError() == nil || *app.GetError() == "") { // DELETING
		return apperrors.NewConcurrentOperationInProgressError("delete operation is in progress")
	}
	// Note: This will be needed when there is async UPDATE supported
	// if app.DeletedAt.IsZero() && app.UpdatedAt.After(app.CreatedAt) && !app.Ready && *app.Error == "" { // UPDATING
	// 	return nil, apperrors.NewInvalidData	Error("another operation is in progress")
	// }

	return nil
}

func (d *directive) prepareRequestData(ctx context.Context, err error, res interface{}) (string, error) {
	tenantID, err := d.tenantLoaderFunc(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve tenant from request")
	}

	app, ok := res.(*graphql.Application)
	if !ok {
		return "", errors.New("entity is not a webhook provider")
	}

	headers, ok := ctx.Value(header.ContextKey).(http.Header)
	if !ok {
		return "", errors.New("failed to retrieve request headers")
	}

	requestData := &RequestData{
		Application: *app,
		TenantID:    tenantID,
		Headers:     headers,
	}

	data, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (d *directive) prepareWebhookIDs(ctx context.Context, err error, operation *Operation, webhookType graphql.WebhookType) ([]string, error) {
	webhooks, err := d.webhookFetcherFunc(ctx, operation.ResourceID)
	if err != nil {
		return nil, err
	}

	webhookIDs := make([]string, 0)
	for _, webhook := range webhooks {
		if graphql.WebhookType(webhook.Type) == webhookType {
			webhookIDs = append(webhookIDs, webhook.ID)
		}
	}
	return webhookIDs, nil
}
