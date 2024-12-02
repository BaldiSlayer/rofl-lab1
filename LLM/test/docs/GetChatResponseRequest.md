# GetChatResponseRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**prompt** | **str** |  | 
**context** | **str** |  | [optional] 
**model** | **str** |  | [optional] [default to 'open-mistral-7b']

## Example

```python
from openapi_client.models.get_chat_response_request import GetChatResponseRequest

# TODO update the JSON string below
json = "{}"
# create an instance of GetChatResponseRequest from a JSON string
get_chat_response_request_instance = GetChatResponseRequest.from_json(json)
# print the JSON string representation of the object
print(GetChatResponseRequest.to_json())

# convert the object into a dict
get_chat_response_request_dict = get_chat_response_request_instance.to_dict()
# create an instance of GetChatResponseRequest from a dict
get_chat_response_request_from_dict = GetChatResponseRequest.from_dict(get_chat_response_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


