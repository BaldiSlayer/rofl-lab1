from app.core.vector_db.faiss_db import convex_indexes


def test_convex_indexes():
    # Тест с обычным случаем
    assert convex_indexes(4, [2, 2, 1]) == 2
    assert convex_indexes(4, [1, 1, 1, 1, 1]) == 4
    assert convex_indexes(2, [1, 1, 1]) == 2
    assert convex_indexes(9, [1]*10) == 9
    assert convex_indexes(0, [1, 1, 1]) == 0

    assert convex_indexes(-1, [1, 2, 3]) == 0
    assert convex_indexes(0, [1]) == 0
    assert convex_indexes(0, [3]) == 0
    assert convex_indexes(0, [10000]) == 0

    assert convex_indexes(1000000, [500000] * 3) == 2
    assert convex_indexes(100, [100] * 10) == 1
    assert convex_indexes(100, [50, 50, 50]) == 2