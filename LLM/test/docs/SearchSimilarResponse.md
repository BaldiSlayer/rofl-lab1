# SearchSimilarResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**result** | [**List[QuestionAnswer]**](QuestionAnswer.md) |  | 

## Example

```python
from openapi_client.models.search_similar_response import SearchSimilarResponse

# TODO update the JSON string below
json = "{}"
# create an instance of SearchSimilarResponse from a JSON string
search_similar_response_instance = SearchSimilarResponse.from_json(json)
# print the JSON string representation of the object
print(SearchSimilarResponse.to_json())

# convert the object into a dict
search_similar_response_dict = search_similar_response_instance.to_dict()
# create an instance of SearchSimilarResponse from a dict
search_similar_response_from_dict = SearchSimilarResponse.from_dict(search_similar_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


