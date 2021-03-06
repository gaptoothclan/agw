package agw

import (
	"net/url"
	"reflect"
	"testing"
)

var ev = ` 
{
    "resource": "/test1/{proxy+}",
    "path": "/test1/test",
    "httpMethod": "POST",
    "headers": {
			"Accept": "*/*",
			"CloudFront-Forwarded-Proto": "https",
			"Content-Type": "application/json"
		},
    "queryStringParameters": {
        "k1": "v1",
        "k2": "v2"
    },
    "pathParameters": {
        "proxy": "test"
    },
    "stageVariables": {
        "lbAlias": "current"
    },
    "requestContext": {
        "path": "/test1/{proxy+}",
        "accountId": "095615327118",
        "resourceId": "ybki7l",
        "stage": "test-invoke-stage",
        "requestId": "test-invoke-request",
        "identity": {
            "cognitoIdentityPoolId": null,
            "cognitoIdentityId": null,
            "apiKey": "test-invoke-api-key",
            "cognitoAuthenticationType": null,
            "userArn": "arn:aws:iam::095615327118:root",
            "apiKeyId": "test-invoke-api-key-id",
            "userAgent": "Apache-HttpClient/4.5.x (Java/1.8.0_144)",
            "accountId": "095615327118",
            "caller": "095615327118",
            "sourceIp": "test-invoke-source-ip",
            "accessKey": "ASIAJTPDCBBJQKRD3FMQ",
            "cognitoAuthenticationProvider": null,
            "user": "095615327118"
        },
        "resourcePath": "/test1/{proxy+}",
        "httpMethod": "POST",
        "apiId": "uorto7w779"
    },
    "body": "{\"key1\":\"value1\"}",
    "isBase64Encoded": false
}
`[1:]

func Test_apiGateParser_queryStringParameters(t *testing.T) {
	var ev2 = `
{
	"queryStringParameters": null
}
`
	uv1 := make(url.Values)
	uv1.Add("k1", "v1")
	uv1.Add("k2", "v2")
	type fields struct {
		content []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   url.Values
	}{
		{"t1", fields{[]byte(ev)}, uv1},
		{"t2", fields{[]byte(ev2)}, make(url.Values)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agp := &APIGateParser{
				content: tt.fields.content,
			}
			if got := agp.QueryStringParameters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiGateParser.queryStringParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiGateParser_url(t *testing.T) {
	type fields struct {
		content []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"t1", fields{[]byte(ev)}, "/test1/test?k1=v1&k2=v2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agp := &APIGateParser{
				content: tt.fields.content,
			}
			if got := agp.Url(); got != tt.want {
				t.Errorf("apiGateParser.url() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiGateParser_body(t *testing.T) {
	type fields struct {
		content []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"t1", fields{[]byte(ev)}, []byte("{\"key1\":\"value1\"}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agp := &APIGateParser{
				content: tt.fields.content,
			}
			if got := agp.Body(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiGateParser.body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIGateParser_StageVariables(t *testing.T) {
	var ev3 = `
{
	"stageVariables": null
},`
	type fields struct {
		content []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{"t1", fields{[]byte(ev1)}, map[string]string{"lbAlias": "current"}},
		{"t2", fields{[]byte(ev3)}, map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agp := &APIGateParser{
				content: tt.fields.content,
			}
			if got := agp.StageVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("APIGateParser.StageVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiGateParser_headers(t *testing.T) {
	name := "t1"
	rawJSON := []byte(ev)

	want := make(map[string]string)
	want["Accept"] = "*/*"
	want["CloudFront-Forwarded-Proto"] = "https"
	want["Content-Type"] = "application/json"

	t.Run(name, func(t *testing.T) {
		agp := &APIGateParser{
			content: rawJSON,
		}
		if got := agp.Headers(); !reflect.DeepEqual(got, want) {
			t.Errorf("apiGateParser.body() = %v, want %v", got, want)
		}
	})
}
