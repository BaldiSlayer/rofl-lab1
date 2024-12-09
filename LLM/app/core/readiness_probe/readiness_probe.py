class ReadinessProbe:
    _value = False

    @classmethod
    def get_value(cls):
        return cls._value

    @classmethod
    def set_value(cls, value: bool):
        cls._value = value
