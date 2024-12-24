# SearchSimilarRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**question** | **str** |  | 

## Example

```python
from openapi_client.models.search_similar_request import SearchSimilarRequest

# TODO update the JSON string below
json = "{}"
# create an instance of SearchSimilarRequest from a JSON string
search_similar_request_instance = SearchSimilarRequest.from_json(json)
# print the JSON string representation of the object
print(SearchSimilarRequest.to_json())

# convert the object into a dict
search_similar_request_dict = search_similar_request_instance.to_dict()
# create an instance of SearchSimilarRequest from a dict
search_similar_request_from_dict = SearchSimilarRequest.from_dict(search_similar_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


