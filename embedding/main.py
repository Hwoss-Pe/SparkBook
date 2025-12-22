from fastapi import FastAPI
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from FlagEmbedding import FlagModel
import uvicorn
from typing import List

app = FastAPI()
model = FlagModel(
    'BAAI/bge-small-zh-v1.5',
    # query_instruction_for_retrieval='为这个句子生成表示以用于检索相关文档：',
    # use_fp16=False,
    device='gpu'
)


class EmbedReq(BaseModel):
    texts: List[str]
    normalize: bool = True
    is_query: bool = True


@app.post("/embed")
def embed(req: EmbedReq):
    texts = [t for t in req.texts if isinstance(t, str) and t.strip()]
    if not texts:
        return {"vectors": [], "dims": 0}
    try:
        if req.is_query:
            vecs = model.encode_queries(texts)
        else:
            vecs = model.encode_corpus(texts)
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})
    if vecs is None or (hasattr(vecs, "__len__") and len(vecs) == 0):
        return {"vectors": [], "dims": 0}
    # print(vecs)
    out = [list(map(float, v)) for v in vecs]
    if req.normalize:
        for i in range(len(out)):
            s = 0.0
            for x in out[i]:
                s += x * x
            if s > 0:
                import math
                n = math.sqrt(s)
                out[i] = [x / n for x in out[i]]
    return {"vectors": out, "dims": len(out[0])}


@app.get("/health")
def health():
    return {"status": "ok"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8088)
