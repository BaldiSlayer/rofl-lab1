# openapi_client.QuestionsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_get_chat_response_get_chat_response_post**](QuestionsApi.md#api_get_chat_response_get_chat_response_post) | **POST** /get_chat_response | Api Get Chat Response
[**api_process_questions_process_questions_post**](QuestionsApi.md#api_process_questions_process_questions_post) | **POST** /search_similar | Api Process Questions


# **api_get_chat_response_get_chat_response_post**
> GetChatResponseResponse api_get_chat_response_get_chat_response_post(get_chat_response_request)

Api Get Chat Response

### Example


```python
import openapi_client
from openapi_client.models.get_chat_response_request import GetChatResponseRequest
from openapi_client.models.get_chat_response_response import GetChatResponseResponse
from openapi_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = openapi_client.QuestionsApi(api_client)
    get_chat_response_request = openapi_client.GetChatResponseRequest() # GetChatResponseRequest | 

    try:
        # Api Get Chat Response
        api_response = api_instance.api_get_chat_response_get_chat_response_post(get_chat_response_request)
        print("The response of QuestionsApi->api_get_chat_response_get_chat_response_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling QuestionsApi->api_get_chat_response_get_chat_response_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **get_chat_response_request** | [**GetChatResponseRequest**](GetChatResponseRequest.md)|  | 

### Return type

[**GetChatResponseResponse**](GetChatResponseResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successful Response |  -  |
**422** | Validation Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_process_questions_process_questions_post**
> SearchSimilarResponse api_process_questions_process_questions_post(search_similar_request)

Api Process Questions

### Example


```python
import openapi_client
from openapi_client.models.search_similar_request import SearchSimilarRequest
from openapi_client.models.search_similar_response import SearchSimilarResponse
from openapi_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost
# See configuration.py for a list of all supported configuration parameters.
configuration = openapi_client.Configuration(
    host = "http://localhost"
)


# Enter a context with an instance of the API client
with openapi_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = openapi_client.QuestionsApi(api_client)
    search_similar_request = openapi_client.SearchSimilarRequest() # SearchSimilarRequest | 

    try:
        # Api Process Questions
        api_response = api_instance.api_process_questions_process_questions_post(search_similar_request)
        print("The response of QuestionsApi->api_process_questions_process_questions_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling QuestionsApi->api_process_questions_process_questions_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **search_similar_request** | [**SearchSimilarRequest**](SearchSimilarRequest.md)|  | 

### Return type

[**SearchSimilarResponse**](SearchSimilarResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successful Response |  -  |
**422** | Validation Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

