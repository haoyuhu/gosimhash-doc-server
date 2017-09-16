# GoSimhash Doc Server for Chinese Documents

[![Build Status](https://travis-ci.org/HaoyuHu/gosimhash-doc-server.svg?branch=master)](https://travis-ci.org/HaoyuHu/gosimhash-doc-server) 
[![License](https://img.shields.io/badge/license-MIT-yellow.svg?style=flat)](http://mit-license.huhaoyu.com)

## Introduction
**GoSimhash Doc Server** has created a backend service with the ability to calculate simhash, compare docs and store simhash fingerprints, providing APIs to serve the **deduplication of Chinese documents**. Inside the server, GoSimhash Doc Server use **Redis** and 
[**Pigeonhole Principle**](https://en.wikipedia.org/wiki/Pigeonhole_principle) to speed up the document similarity matching.

## Build
```
// fetch gosimhash-doc-server
go get github.com/HaoyuHu/gosimhash-doc-server
// build it in the root directory of gosimhash-doc-server project
go build
```

## Configuration
### dict/*.utf8
You can put custom dicts in `dict/`, such as hmm model dict, idf dict, user custom dict, stop words. The project provide default dicts.

### config/common.json
This is **global configurations** of gosimhash-doc-server.

```json
{
  "hash_type": 1,
  "hmm_dict": "dict/hmm_model.utf8",
  "idf_dict": "dict/idf.utf8",
  "user_dict": "dict/jieba.dict.utf8",
  "stop_words": "dict/stop_words.utf8",
  "simhash_limit": 3,
  "host": "",
  "port": "8000"
}
```

* **hash_type**: gosimhash-doc-server provide two types of inner hash method for simhash, [**SIPHASH**](https://en.wikipedia.org/wiki/SipHash)(**0**) and [**JENKINS**](https://en.wikipedia.org/wiki/Jenkins_hash_function)(**1**, default value);
* **hmm_dict**: [HMM model](https://en.wikipedia.org/wiki/Hidden_Markov_model) for [gojieba](https://github.com/yanyiwu/gojieba);
* **idf_dict**: [idf dict](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) for [gojieba](https://github.com/yanyiwu/gojieba), which is used for calculating the weight of words;
* **user_dict**: custom dict for [gojieba](https://github.com/yanyiwu/gojieba);
* **stop_words**: custom stop word list for [gojieba](https://github.com/yanyiwu/gojieba);
* **simhash_limit**: for determining wether documents is similar, default value is **3**, you can choose **1/3/7/15**;
* **host**: server host, GoSimhash Doc Server will find **HOST** in **Environment variables** at first;
* **port**: server port, GoSimhash Doc Server will find **PORT** in **Environment variables** at first;

### config/redis.json
```json
{
  "host": "localhost",
  "port": 6379,
  "passwd": ""
}
```

* **host**: Redis host, default `"localhost"`;
* **port**: Redis port, default `6379`;
* **passwd**: Redis password, default `""`;

## Usage
After building project, a executable file named `gosimhash-doc-server` will be in the root directory of the project.
After configuring common.json and redis.json, the directory structure is like below:
```
// you can remove all other useless files
+-- config
|   +-- common.json
|   +-- redis.json
+-- dict
|   +-- hmm_model.utf8
|   +-- idf.utf8
|   +-- jieba.dict.utf8
|   +-- stop_words.utf8
+-- gosimhash-doc-server // executable file
```

## Startup
Make sure your Redis service is running on the correct host and port.

```
./gosimhash-doc-server
```

## APIs

### POST /simhash
Calculate simhash fingerprint of the current document. The fingerprint will **NOT** insert into the Redis for calculating the similarity of documents.

#### Post Params:

* **doc**: Document content;
* **top_n**: Use the top N words with the largest weight in the word segmentation of the current document to calculate the simhash fingerprint.

#### Json Response:
on success:
```json
{
  "code": 0,
  "err_message": "ok",
  "data": 7422112736388306927
}
```
on error:
```json
{
  "code": -1,
  "err_message": "Empty doc or top_n",
  "data": null
}
```

### POST /distance
Calculate the distance of two documents. The fingerprint will **NOT** insert into the Redis for calculating the similarity of documents.

#### Post Params:

* **first_doc**: First document content;
* **second_doc**: Second document content;
* **top_n**: Use the top N words with the largest weight in the word segmentation of the current document to calculate the simhash fingerprint.

#### Json Response:
on success:
```json
{
  "code": 0,
  "err_message": "ok",
  "data": {
    "distance": 23,
    "first_simhash": 7422112736388306927,
    "second_simhash": 5700964324719740863
    }
}
```
on error:
```json
{
  "code": -1,
  "err_message": "Empty doc or top_n",
  "data": null
}
```

### POST /identify
Calculate simhash for current document. The fingerprint will insert into the Redis for calculating the similarity of documents if there is no similar documents in Redis (simhash distance <= simhash_limit).

#### Post Params:

* **doc_id**: Document id;
* **doc**: Document content;
* **age**: Age of current document. Document fingerprint will be permanently present in the fingerprint library if this parameter is not passed;
* **top_n**: Use the top N words with the largest weight in the word segmentation of the current document to calculate the simhash fingerprint.

#### Json Response:
on success(no similar document):

```json
{
  "code": 0,
  "err_message": "ok",
  "data": {
  "has_similar_doc": false
  }
}
```

on success(has similar document):
```json
{
  "code": 0,
  "err_message": "ok",
  "data": {
  "has_similar_doc": true,
  "similar_doc_id": "12345678"
  }
}
```

on error:
```json
{
  "code": -1,
  "err_message": "Empty doc or top_n",
  "data": null
}
```
