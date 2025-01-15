# ProcessQuestionsRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**questions_list** | [**List[QuestionAnswer]**](QuestionAnswer.md) |  |
**use_saved** | **bool** |  | [optional] [default to False]
**filename** | **str** |  | [optional] [default to 'vectorized_data']

## Example

```python
from llm_client.models.process_questions_request import ProcessQuestionsRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ProcessQuestionsRequest from a JSON string
process_questions_request_instance = ProcessQuestionsRequest.from_json(json)
# print the JSON string representation of the object
print(ProcessQuestionsRequest.to_json())

# convert the object into a dict
process_questions_request_dict = process_questions_request_instance.to_dict()
# create an instance of ProcessQuestionsRequest from a dict
process_questions_request_from_dict = ProcessQuestionsRequest.from_dict(process_questions_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
