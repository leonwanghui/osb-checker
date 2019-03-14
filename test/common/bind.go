package common

import (
	"testing"

	v2 "github.com/openservicebrokerapi/osb-checker/autogenerated/models"
	apiclient "github.com/openservicebrokerapi/osb-checker/client"
	. "github.com/openservicebrokerapi/osb-checker/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBind(
	t *testing.T,
	instanceID, bindingID string,
	req *v2.ServiceBindingRequest,
	async bool,
) {
	t.Parallel()

	Convey("BINDING - request syntax", t, func() {

		So(testAPIVersionHeader(GenerateBindingURL(instanceID, bindingID), "PUT"), ShouldEqual, nil)
		So(testAuthentication(GenerateBindingURL(instanceID, bindingID), "PUT"), ShouldEqual, nil)
		So(testAsyncParameters(GenerateBindingURL(instanceID, bindingID), "PUT"), ShouldEqual, nil)

		var tempBody = new(v2.ServiceBindingRequest)
		Convey("should reject if missing service_id", func() {
			deepCopy(req, tempBody)
			*tempBody.ServiceID = ""
			code, _, err := apiclient.Default.Bind(instanceID, bindingID, tempBody, async)

			So(err, ShouldEqual, nil)
			So(code, ShouldEqual, 400)
		})

		Convey("should reject if missing plan_id", func() {
			deepCopy(req, tempBody)
			*tempBody.PlanID = ""
			code, _, err := apiclient.Default.Bind(instanceID, bindingID, tempBody, async)

			So(err, ShouldEqual, nil)
			So(code, ShouldEqual, 400)
		})

		Convey("should reject if service_id is invalid", func() {
			deepCopy(req, tempBody)
			*tempBody.ServiceID = "xxxx-xxxx"
			code, _, err := apiclient.Default.Bind(instanceID, bindingID, tempBody, async)

			So(err, ShouldEqual, nil)
			So(code, ShouldEqual, 400)
		})

		Convey("should reject if paln_id is invalid", func() {
			deepCopy(req, tempBody)
			*tempBody.PlanID = "xxxx-xxxx"
			code, _, err := apiclient.Default.Bind(instanceID, bindingID, tempBody, async)

			So(err, ShouldEqual, nil)
			So(code, ShouldEqual, 400)
		})

		Convey("should accept a valid binding request", func() {
			deepCopy(req, tempBody)
			code, body, err := apiclient.Default.Bind(instanceID, bindingID, tempBody, async)

			So(err, ShouldEqual, nil)
			if async {
				So(code, ShouldEqual, 202)
			} else {
				So(code, ShouldEqual, 201)
			}
			So(testJSONSchema(body), ShouldEqual, nil)
		})

		Convey("should return 200 OK when binding Id with same instance Id exists with identical properties", func() {
			deepCopy(req, tempBody)
			code, body, err := apiclient.Default.Bind(instanceID, bindingID, tempBody, async)

			So(err, ShouldEqual, nil)
			So(code, ShouldEqual, 200)
			So(testJSONSchema(body), ShouldEqual, nil)
		})
	})
}
