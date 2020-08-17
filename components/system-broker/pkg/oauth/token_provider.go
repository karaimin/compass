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

package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	httputils "github.com/kyma-incubator/compass/components/system-broker/pkg/http"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	contentTypeHeader                = "Content-Type"
	contentTypeApplicationURLEncoded = "application/x-www-form-urlencoded"

	grantTypeFieldName   = "grant_type"
	credentialsGrantType = "client_credentials"

	scopeFieldName = "scope"
	scopes         = "runtime:read runtime:write"

	clientIDKey       = "client_id"
	clientSecretKey   = "client_secret"
	tokensEndpointKey = "tokens_endpoint"
)

type RequestProvider interface {
	Provide(ctx context.Context, input httputils.RequestInput) (*http.Request, error)
}

type OAuthTokenProvider struct {
	httpClient      *http.Client
	requestProvider RequestProvider
	k8sClient       client.Client
	secretName      string
	secretNamespace string
}

type credentials struct {
	clientID       string
	clientSecret   string
	tokensEndpoint string
}

func NewOAuthTokenProviderFromSecret(config *Config, httpClient *http.Client, requestProvider RequestProvider, k8sClient client.Client) *OAuthTokenProvider {
	return &OAuthTokenProvider{
		httpClient:      httpClient,
		requestProvider: requestProvider,
		k8sClient:       k8sClient,
		secretName:      config.SecretName,
		secretNamespace: config.SecretNamespace,
	}
}

func (c *OAuthTokenProvider) GetAuthorizationToken(ctx context.Context) (httputils.Token, error) {
	credentials, err := c.extractOAuthClientFromSecret(ctx)
	if err != nil {
		return httputils.Token{}, errors.Wrap(err, "while get credentials from secret")
	}

	return c.getAuthorizationToken(ctx, credentials)
}

func (c *OAuthTokenProvider) WaitForCredentials(ctx context.Context) error {
	err := wait.Poll(time.Second, time.Minute*3, func() (bool, error) {
		secret := &v1.Secret{}
		err := c.k8sClient.Get(ctx, client.ObjectKey{
			Namespace: c.secretNamespace,
			Name:      c.secretName,
		}, secret)
		// it fails on connection-refused error on first call and it restarts our application.
		if err != nil {
			log.C(ctx).Warnf("secret %s not found", c.secretName)
			return false, nil
		}
		return true, nil
	})

	return errors.Wrapf(err, "while waiting for secret %s", c.secretName)
}

func (c *OAuthTokenProvider) extractOAuthClientFromSecret(ctx context.Context) (credentials, error) {
	secret := &v1.Secret{}
	err := c.k8sClient.Get(ctx, client.ObjectKey{
		Namespace: c.secretNamespace,
		Name:      c.secretName,
	}, secret)
	if err != nil {
		return credentials{}, err
	}

	return credentials{
		clientID:       string(secret.Data[clientIDKey]),
		clientSecret:   string(secret.Data[clientSecretKey]),
		tokensEndpoint: string(secret.Data[tokensEndpointKey]),
	}, nil
}

func (c *OAuthTokenProvider) getAuthorizationToken(ctx context.Context, credentials credentials) (httputils.Token, error) {
	log.C(ctx).Infof("Getting authorization token from endpoint: %s", credentials.tokensEndpoint)

	form := url.Values{}
	form.Add(grantTypeFieldName, credentialsGrantType)
	form.Add(scopeFieldName, scopes)
	body := strings.NewReader(form.Encode())

	headers := map[string]string{
		contentTypeHeader: contentTypeApplicationURLEncoded,
	}

	request, err := c.requestProvider.Provide(ctx, httputils.RequestInput{
		Method:  http.MethodPost,
		URL:     credentials.tokensEndpoint,
		Body:    body,
		Headers: headers,
	})
	if err != nil {
		return httputils.Token{}, errors.Wrap(err, "while creating authorisation token request")
	}

	request.SetBasicAuth(credentials.clientID, credentials.clientSecret)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return httputils.Token{}, errors.Wrap(err, "while send request to token endpoint")
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.C(ctx).Warn("Cannot close connection body inside oauth client")
		}
	}()

	if response.StatusCode != http.StatusOK {
		return httputils.Token{}, fmt.Errorf("while calling to token endpoint: unexpected status code, %d, %s", response.StatusCode, response.Status)
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return httputils.Token{}, errors.Wrapf(err, "while reading token response body from %q", credentials.tokensEndpoint)
	}

	tokenResponse := httputils.Token{}
	err = json.Unmarshal(respBody, &tokenResponse)
	if err != nil {
		return httputils.Token{}, errors.Wrap(err, "while unmarshalling token response body")
	}

	if tokenResponse.AccessToken == "" {
		return httputils.Token{}, errors.New("while fetching token: access token from oauth client is empty")
	}

	log.C(ctx).Info("Successfully unmarshal response oauth token for accessing Director")
	tokenResponse.Expiration += time.Now().Unix()

	return tokenResponse, nil
}
