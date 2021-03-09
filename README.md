# lector
Charactar normalisation service that renders unicode confusables and send back the string via ocr and makes a judgement about profanity.

There is [here](Lector.postman_collection.json) a sample Postman collection.

To test it on your local machine just forward the service to your localhost and try the examples.

Current integration environment access is needed:

```bash
kubectl --context i20210203.int.freeletics.com -n lector port-forward service/lector 8000:8000
```

Sample payload:

```json
{"toCheck": "ꜰᴜᴄᴋ ᴍᴇ"}
```

Sample Response

```json
{
    "ocr": {
        "string": "FUCK ME",
        "profan": true
    },
    "raw": {
        "string": "ꜰᴜᴄᴋ ᴍᴇ",
        "profan": false
    },
    "transcribed": {
        "string": "ꜰucĸ ʍᴇ",
        "profan": false
    }
}
```



Response struct:

```go
type Response struct {
	Ocr struct {
		String string `json:"string"`
		Profan bool   `json:"profan"`
	} `json:"ocr"`
	Raw struct {
		String string `json:"string"`
		Profan bool   `json:"profan"`
	} `json:"raw"`
	Transcribed struct {
		String string `json:"string"`
		Profan bool   `json:"profan"`
	} `json:"transcribed"`
}
```



Confusbales in unicode are characters that look a like another one.

http://www.unicode.org/Public/security/latest/confusables.txt

If you like to try more sophisticated strings you can create one on your own [here](https://www.irongeek.com/homoglyph-attack-generator.php)

One possible answer would be this lector service.



Credits to:



![](img/maxresdefault.jpg)