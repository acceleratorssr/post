from elasticsearch import Elasticsearch

ARTICLE_INDEX_NAME = "article_index"

def init_es_client():
    es = Elasticsearch(
        ["http://localhost:9200"],
    )
    return es

def delete_index(es):
    if es.indices.exists(index=ARTICLE_INDEX_NAME):
        es.indices.delete(index=ARTICLE_INDEX_NAME)
        print("Index deleted")
    else:
        print(f"Index '{ARTICLE_INDEX_NAME}' does not exist.")

if __name__ == "__main__":
    es_client = init_es_client()
    delete_index(es_client)
