import sys
import math
import requests


def norm(v):
    return math.sqrt(sum(x * x for x in v))


def call(texts, is_query=True, url="http://localhost:8088/embed"):
    r = requests.post(url, json={"texts": texts, "normalize": True, "is_query": is_query}, timeout=5)
    r.raise_for_status()
    data = r.json()
    vecs = data["vectors"]
    dims = len(vecs[0]) if vecs else 0
    return dims, vecs


def main():
    q = sys.argv[1] if len(sys.argv) > 1 else "向量检索测试"
    dims, vecs = call([q], is_query=True)
    print("query dims:", dims)
    print("query norm:", round(norm(vecs[0]), 6))
    dims2, vecs2 = call(["这是一个用于检索的文档内容示例"], is_query=False)
    print("corpus dims:", dims2)
    print("corpus norm:", round(norm(vecs2[0]), 6))
    dot = sum(vecs[0][i] * vecs2[0][i] for i in range(min(len(vecs[0]), len(vecs2[0]))))
    print("cosine approx:", round(dot, 6))


if __name__ == "__main__":
    main()
