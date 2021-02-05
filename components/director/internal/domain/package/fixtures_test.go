package mp_package_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	mp_package "github.com/kyma-incubator/compass/components/director/internal/domain/package"
	"github.com/kyma-incubator/compass/components/director/internal/model"
)

const (
	packageID        = "packageID"
	tenantID         = "tenantID"
	appID            = "appID"
	ordID            = "com.compass.v1"
	externalTenantID = "externalTenantID"
)

func fixEntityPackage() *mp_package.Entity {
	return &mp_package.Entity{
		ID:                packageID,
		TenantID:          tenantID,
		ApplicationID:     appID,
		OrdID:             ordID,
		Vendor:            sql.NullString{String: "vendorID", Valid: true},
		Title:             "title",
		ShortDescription:  "short desc",
		Description:       "desc",
		Version:           "v1.0.5",
		PackageLinks:      json.RawMessage("{}"),
		Links:             json.RawMessage("[]"),
		LicenceType:       sql.NullString{String: "test", Valid: true},
		Tags:              json.RawMessage("[]"),
		Countries:         json.RawMessage("[]"),
		Labels:            json.RawMessage("{}"),
		PolicyLevel:       "test",
		CustomPolicyLevel: sql.NullString{},
		PartOfProducts:    json.RawMessage("[\"test\"]"),
		LineOfBusiness:    json.RawMessage("[]"),
		Industry:          json.RawMessage("[]"),
	}
}

func fixPackageModel() *model.Package {
	vendorID := "vendorID"
	licenceType := "test"
	return &model.Package{
		ID:                packageID,
		TenantID:          tenantID,
		ApplicationID:     appID,
		OrdID:             ordID,
		Vendor:            &vendorID,
		Title:             "title",
		ShortDescription:  "short desc",
		Description:       "desc",
		Version:           "v1.0.5",
		PackageLinks:      json.RawMessage("{}"),
		Links:             json.RawMessage("[]"),
		LicenceType:       &licenceType,
		Tags:              json.RawMessage("[]"),
		Countries:         json.RawMessage("[]"),
		Labels:            json.RawMessage("{}"),
		PolicyLevel:       "test",
		CustomPolicyLevel: nil,
		PartOfProducts:    json.RawMessage("[\"test\"]"),
		LineOfBusiness:    json.RawMessage("[]"),
		Industry:          json.RawMessage("[]"),
	}
}

func fixPackageModelInput() *model.PackageInput {
	vendorID := "vendorID"
	licenceType := "test"
	return &model.PackageInput{
		OrdID:             ordID,
		Vendor:            &vendorID,
		Title:             "title",
		ShortDescription:  "short desc",
		Description:       "desc",
		Version:           "v1.0.5",
		PackageLinks:      json.RawMessage("{}"),
		Links:             json.RawMessage("[]"),
		LicenceType:       &licenceType,
		Tags:              json.RawMessage("[]"),
		Countries:         json.RawMessage("[]"),
		Labels:            json.RawMessage("{}"),
		PolicyLevel:       "test",
		CustomPolicyLevel: nil,
		PartOfProducts:    json.RawMessage("[\"test\"]"),
		LineOfBusiness:    json.RawMessage("[]"),
		Industry:          json.RawMessage("[]"),
	}
}

func fixPackageColumns() []string {
	return []string{"id", "tenant_id", "app_id", "ord_id", "vendor", "title", "short_description",
		"description", "version", "package_links", "links", "licence_type", "tags", "countries", "labels", "policy_level",
		"custom_policy_level", "part_of_products", "line_of_business", "industry"}
}

func fixPackageRow() []driver.Value {
	return []driver.Value{packageID, tenantID, appID, ordID, "vendorID", "title", "short desc", "desc", "v1.0.5",
		json.RawMessage("{}"), json.RawMessage("[]"), "test", json.RawMessage("[]"), json.RawMessage("[]"), json.RawMessage("{}"),
		"test", nil, json.RawMessage("[\"test\"]"), json.RawMessage("[]"), json.RawMessage("[]")}
}

func fixPackageUpdateArgs() []driver.Value {
	return []driver.Value{"vendorID", "title", "short desc", "desc", "v1.0.5", json.RawMessage("{}"), json.RawMessage("[]"),
		"test", json.RawMessage("[]"), json.RawMessage("[]"), json.RawMessage("{}"), "test", nil, json.RawMessage("[\"test\"]"), json.RawMessage("[]"), json.RawMessage("[]")}
}
