# QuestionAnswer


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**question** | **str** |  | 
**answer** | **str** |  | 

## Example

```python
from openapi_client.models.question_answer import QuestionAnswer

# TODO update the JSON string below
json = "{}"
# create an instance of QuestionAnswer from a JSON string
question_answer_instance = QuestionAnswer.from_json(json)
# print the JSON string representation of the object
print(QuestionAnswer.to_json())

# convert the object into a dict
question_answer_dict = question_answer_instance.to_dict()
# create an instance of QuestionAnswer from a dict
question_answer_from_dict = QuestionAnswer.from_dict(question_answer_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


