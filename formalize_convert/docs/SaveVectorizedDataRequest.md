# SaveVectorizedDataRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**data** | [**List[QuestionAnswer]**](QuestionAnswer.md) |  |
**embeddings** | **List[List[float]]** |  |
**filename** | **str** |  | [optional] [default to 'vectorized_data']

## Example

```python
from llm_client.models.save_vectorized_data_request import SaveVectorizedDataRequest

# TODO update the JSON string below
json = "{}"
# create an instance of SaveVectorizedDataRequest from a JSON string
save_vectorized_data_request_instance = SaveVectorizedDataRequest.from_json(json)
# print the JSON string representation of the object
print(SaveVectorizedDataRequest.to_json())

# convert the object into a dict
save_vectorized_data_request_dict = save_vectorized_data_request_instance.to_dict()
# create an instance of SaveVectorizedDataRequest from a dict
save_vectorized_data_request_from_dict = SaveVectorizedDataRequest.from_dict(save_vectorized_data_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
