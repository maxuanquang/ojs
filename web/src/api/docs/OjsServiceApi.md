# OjsProto.OjsServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ojsServiceCreateAccount**](OjsServiceApi.md#ojsServiceCreateAccount) | **POST** /api/v1/accounts | 
[**ojsServiceCreateProblem**](OjsServiceApi.md#ojsServiceCreateProblem) | **POST** /api/v1/problems | 
[**ojsServiceCreateSession**](OjsServiceApi.md#ojsServiceCreateSession) | **POST** /api/v1/sessions | 
[**ojsServiceCreateSubmission**](OjsServiceApi.md#ojsServiceCreateSubmission) | **POST** /api/v1/submissions | 
[**ojsServiceCreateTestCase**](OjsServiceApi.md#ojsServiceCreateTestCase) | **POST** /api/v1/test-cases | 
[**ojsServiceDeleteProblem**](OjsServiceApi.md#ojsServiceDeleteProblem) | **DELETE** /api/v1/problems/{id} | 
[**ojsServiceDeleteSession**](OjsServiceApi.md#ojsServiceDeleteSession) | **DELETE** /api/v1/sessions | 
[**ojsServiceDeleteTestCase**](OjsServiceApi.md#ojsServiceDeleteTestCase) | **DELETE** /api/v1/test-cases/{id} | 
[**ojsServiceGetAccount**](OjsServiceApi.md#ojsServiceGetAccount) | **GET** /api/v1/accounts/{id} | 
[**ojsServiceGetAccountProblemSubmissionList**](OjsServiceApi.md#ojsServiceGetAccountProblemSubmissionList) | **GET** /api/v1/accounts/{accountId}/problems/{problemId}/submissions | 
[**ojsServiceGetProblem**](OjsServiceApi.md#ojsServiceGetProblem) | **GET** /api/v1/problems/{id} | 
[**ojsServiceGetProblemList**](OjsServiceApi.md#ojsServiceGetProblemList) | **GET** /api/v1/problems | 
[**ojsServiceGetProblemSubmissionList**](OjsServiceApi.md#ojsServiceGetProblemSubmissionList) | **GET** /api/v1/problems/{id}/submissions | 
[**ojsServiceGetProblemTestCaseList**](OjsServiceApi.md#ojsServiceGetProblemTestCaseList) | **GET** /api/v1/problems/{id}/test-cases | 
[**ojsServiceGetServerInfo**](OjsServiceApi.md#ojsServiceGetServerInfo) | **GET** /api/v1/info | 
[**ojsServiceGetSubmission**](OjsServiceApi.md#ojsServiceGetSubmission) | **GET** /api/v1/submissions/{id} | 
[**ojsServiceGetSubmissionList**](OjsServiceApi.md#ojsServiceGetSubmissionList) | **GET** /api/v1/submissions | 
[**ojsServiceGetTestCase**](OjsServiceApi.md#ojsServiceGetTestCase) | **GET** /api/v1/test-cases/{id} | 
[**ojsServiceUpdateProblem**](OjsServiceApi.md#ojsServiceUpdateProblem) | **PUT** /api/v1/problems/{id} | 
[**ojsServiceUpdateTestCase**](OjsServiceApi.md#ojsServiceUpdateTestCase) | **PUT** /api/v1/test-cases/{id} | 



## ojsServiceCreateAccount

> OjsCreateAccountResponse ojsServiceCreateAccount(body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let body = new OjsProto.OjsCreateAccountRequest(); // OjsCreateAccountRequest | 
apiInstance.ojsServiceCreateAccount(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OjsCreateAccountRequest**](OjsCreateAccountRequest.md)|  | 

### Return type

[**OjsCreateAccountResponse**](OjsCreateAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## ojsServiceCreateProblem

> OjsCreateProblemResponse ojsServiceCreateProblem(body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let body = new OjsProto.OjsCreateProblemRequest(); // OjsCreateProblemRequest | 
apiInstance.ojsServiceCreateProblem(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OjsCreateProblemRequest**](OjsCreateProblemRequest.md)|  | 

### Return type

[**OjsCreateProblemResponse**](OjsCreateProblemResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## ojsServiceCreateSession

> OjsCreateSessionResponse ojsServiceCreateSession(body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let body = new OjsProto.OjsCreateSessionRequest(); // OjsCreateSessionRequest | 
apiInstance.ojsServiceCreateSession(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OjsCreateSessionRequest**](OjsCreateSessionRequest.md)|  | 

### Return type

[**OjsCreateSessionResponse**](OjsCreateSessionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## ojsServiceCreateSubmission

> OjsCreateSubmissionResponse ojsServiceCreateSubmission(body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let body = new OjsProto.OjsCreateSubmissionRequest(); // OjsCreateSubmissionRequest | 
apiInstance.ojsServiceCreateSubmission(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OjsCreateSubmissionRequest**](OjsCreateSubmissionRequest.md)|  | 

### Return type

[**OjsCreateSubmissionResponse**](OjsCreateSubmissionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## ojsServiceCreateTestCase

> OjsCreateTestCaseResponse ojsServiceCreateTestCase(body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let body = new OjsProto.OjsCreateTestCaseRequest(); // OjsCreateTestCaseRequest | 
apiInstance.ojsServiceCreateTestCase(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**OjsCreateTestCaseRequest**](OjsCreateTestCaseRequest.md)|  | 

### Return type

[**OjsCreateTestCaseResponse**](OjsCreateTestCaseResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## ojsServiceDeleteProblem

> Object ojsServiceDeleteProblem(id)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
apiInstance.ojsServiceDeleteProblem(id, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceDeleteSession

> Object ojsServiceDeleteSession()



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
apiInstance.ojsServiceDeleteSession((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceDeleteTestCase

> Object ojsServiceDeleteTestCase(id)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
apiInstance.ojsServiceDeleteTestCase(id, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetAccount

> OjsGetAccountResponse ojsServiceGetAccount(id)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
apiInstance.ojsServiceGetAccount(id, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 

### Return type

[**OjsGetAccountResponse**](OjsGetAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetAccountProblemSubmissionList

> OjsGetAccountProblemSubmissionListResponse ojsServiceGetAccountProblemSubmissionList(accountId, problemId, opts)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let accountId = "accountId_example"; // String | 
let problemId = "problemId_example"; // String | 
let opts = {
  'offset': "offset_example", // String | 
  'limit': "limit_example" // String | 
};
apiInstance.ojsServiceGetAccountProblemSubmissionList(accountId, problemId, opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **String**|  | 
 **problemId** | **String**|  | 
 **offset** | **String**|  | [optional] 
 **limit** | **String**|  | [optional] 

### Return type

[**OjsGetAccountProblemSubmissionListResponse**](OjsGetAccountProblemSubmissionListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetProblem

> OjsGetProblemResponse ojsServiceGetProblem(id)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
apiInstance.ojsServiceGetProblem(id, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 

### Return type

[**OjsGetProblemResponse**](OjsGetProblemResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetProblemList

> OjsGetProblemListResponse ojsServiceGetProblemList(opts)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let opts = {
  'offset': "offset_example", // String | 
  'limit': "limit_example" // String | 
};
apiInstance.ojsServiceGetProblemList(opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **String**|  | [optional] 
 **limit** | **String**|  | [optional] 

### Return type

[**OjsGetProblemListResponse**](OjsGetProblemListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetProblemSubmissionList

> OjsGetProblemSubmissionListResponse ojsServiceGetProblemSubmissionList(id, opts)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
let opts = {
  'offset': "offset_example", // String | 
  'limit': "limit_example" // String | 
};
apiInstance.ojsServiceGetProblemSubmissionList(id, opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **offset** | **String**|  | [optional] 
 **limit** | **String**|  | [optional] 

### Return type

[**OjsGetProblemSubmissionListResponse**](OjsGetProblemSubmissionListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetProblemTestCaseList

> OjsGetProblemTestCaseListResponse ojsServiceGetProblemTestCaseList(id, opts)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
let opts = {
  'offset': "offset_example", // String | 
  'limit': "limit_example" // String | 
};
apiInstance.ojsServiceGetProblemTestCaseList(id, opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **offset** | **String**|  | [optional] 
 **limit** | **String**|  | [optional] 

### Return type

[**OjsGetProblemTestCaseListResponse**](OjsGetProblemTestCaseListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetServerInfo

> Object ojsServiceGetServerInfo()



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
apiInstance.ojsServiceGetServerInfo((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetSubmission

> OjsGetSubmissionResponse ojsServiceGetSubmission(id)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
apiInstance.ojsServiceGetSubmission(id, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 

### Return type

[**OjsGetSubmissionResponse**](OjsGetSubmissionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetSubmissionList

> OjsGetSubmissionListResponse ojsServiceGetSubmissionList(opts)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let opts = {
  'offset': "offset_example", // String | 
  'limit': "limit_example" // String | 
};
apiInstance.ojsServiceGetSubmissionList(opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **String**|  | [optional] 
 **limit** | **String**|  | [optional] 

### Return type

[**OjsGetSubmissionListResponse**](OjsGetSubmissionListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceGetTestCase

> OjsGetTestCaseResponse ojsServiceGetTestCase(id)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
apiInstance.ojsServiceGetTestCase(id, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 

### Return type

[**OjsGetTestCaseResponse**](OjsGetTestCaseResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## ojsServiceUpdateProblem

> OjsUpdateProblemResponse ojsServiceUpdateProblem(id, body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
let body = new OjsProto.OjsServiceUpdateProblemBody(); // OjsServiceUpdateProblemBody | 
apiInstance.ojsServiceUpdateProblem(id, body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **body** | [**OjsServiceUpdateProblemBody**](OjsServiceUpdateProblemBody.md)|  | 

### Return type

[**OjsUpdateProblemResponse**](OjsUpdateProblemResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## ojsServiceUpdateTestCase

> OjsUpdateTestCaseResponse ojsServiceUpdateTestCase(id, body)



### Example

```javascript
import OjsProto from 'ojs_proto';

let apiInstance = new OjsProto.OjsServiceApi();
let id = "id_example"; // String | 
let body = new OjsProto.OjsServiceUpdateTestCaseBody(); // OjsServiceUpdateTestCaseBody | 
apiInstance.ojsServiceUpdateTestCase(id, body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **body** | [**OjsServiceUpdateTestCaseBody**](OjsServiceUpdateTestCaseBody.md)|  | 

### Return type

[**OjsUpdateTestCaseResponse**](OjsUpdateTestCaseResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

