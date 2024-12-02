# AddQuestionsRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**new_questions** | [**List[QuestionAnswer]**](QuestionAnswer.md) |  |
**filename** | **str** |  | [optional] [default to 'vectorized_data']

## Example

```python
from llm_client.models.add_questions_request import AddQuestionsRequest

# TODO update the JSON string below
json = "{}"
# create an instance of AddQuestionsRequest from a JSON string
add_questions_request_instance = AddQuestionsRequest.from_json(json)
# print the JSON string representation of the object
print(AddQuestionsRequest.to_json())

# convert the object into a dict
add_questions_request_dict = add_questions_request_instance.to_dict()
# create an instance of AddQuestionsRequest from a dict
add_questions_request_from_dict = AddQuestionsRequest.from_dict(add_questions_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
