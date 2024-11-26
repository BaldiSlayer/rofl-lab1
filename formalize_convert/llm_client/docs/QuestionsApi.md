# openapi_client.QuestionsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_add_questions_add_questions_post**](QuestionsApi.md#api_add_questions_add_questions_post) | **POST** /add_questions | Api Add Questions
[**api_get_chat_response_get_chat_response_post**](QuestionsApi.md#api_get_chat_response_get_chat_response_post) | **POST** /get_chat_response | Api Get Chat Response
[**api_process_questions_process_questions_post**](QuestionsApi.md#api_process_questions_process_questions_post) | **POST** /process_questions | Api Process Questions
[**api_save_vectorized_data_save_vectorized_data_post**](QuestionsApi.md#api_save_vectorized_data_save_vectorized_data_post) | **POST** /save_vectorized_data | Api Save Vectorized Data


# **api_add_questions_add_questions_post**
> object api_add_questions_add_questions_post(add_questions_request)

Api Add Questions

### Example


```python
import openapi_client
from openapi_client.models.add_questions_request import AddQuestionsRequest
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
    add_questions_request = openapi_client.AddQuestionsRequest() # AddQuestionsRequest |

    try:
        # Api Add Questions
        api_response = api_instance.api_add_questions_add_questions_post(add_questions_request)
        print("The response of QuestionsApi->api_add_questions_add_questions_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling QuestionsApi->api_add_questions_add_questions_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **add_questions_request** | [**AddQuestionsRequest**](AddQuestionsRequest.md)|  |

### Return type

**object**

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
> ProcessQuestionsResponse api_process_questions_process_questions_post(process_questions_request)

Api Process Questions

### Example


```python
import openapi_client
from openapi_client.models.process_questions_request import ProcessQuestionsRequest
from openapi_client.models.process_questions_response import ProcessQuestionsResponse
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
    process_questions_request = openapi_client.ProcessQuestionsRequest() # ProcessQuestionsRequest |

    try:
        # Api Process Questions
        api_response = api_instance.api_process_questions_process_questions_post(process_questions_request)
        print("The response of QuestionsApi->api_process_questions_process_questions_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling QuestionsApi->api_process_questions_process_questions_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **process_questions_request** | [**ProcessQuestionsRequest**](ProcessQuestionsRequest.md)|  |

### Return type

[**ProcessQuestionsResponse**](ProcessQuestionsResponse.md)

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

# **api_save_vectorized_data_save_vectorized_data_post**
> object api_save_vectorized_data_save_vectorized_data_post(save_vectorized_data_request)

Api Save Vectorized Data

### Example


```python
import openapi_client
from openapi_client.models.save_vectorized_data_request import SaveVectorizedDataRequest
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
    save_vectorized_data_request = openapi_client.SaveVectorizedDataRequest() # SaveVectorizedDataRequest |

    try:
        # Api Save Vectorized Data
        api_response = api_instance.api_save_vectorized_data_save_vectorized_data_post(save_vectorized_data_request)
        print("The response of QuestionsApi->api_save_vectorized_data_save_vectorized_data_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling QuestionsApi->api_save_vectorized_data_save_vectorized_data_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **save_vectorized_data_request** | [**SaveVectorizedDataRequest**](SaveVectorizedDataRequest.md)|  |

### Return type

**object**

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
