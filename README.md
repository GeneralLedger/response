<a href="https://github.com/GeneralLedger"><p align="center"><img src="https://user-images.githubusercontent.com/2105067/62828744-96c37a00-bba2-11e9-9c11-ea95f6ab4022.png" alt="General Ledger" width="160px"/></p></a>
<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg" alt="MPL 2.0"></img></a>
</p>

# Response

Response is a no-frills standardized response body wrapper with some added utility to help manage writing to a [http.ResponseWriter](https://golang.org/pkg/net/http/#ResponseWriter).

## Output Interface

```javascript
{
  "status_code": 200,
  "status_text": "OK",
  "error_details": "Invalid email",
  "result": {
    // ...
  }
}
```

## Usage

```go
func MyHandlerFunc(w http.ResponseWriter, r *http.Request) {
    resp := response.New(w)
    defer resp.Output()

    models, err := getModels()
    if err != nil {
        resp.SetResult(http.StatusInternalServerError, nil)
            .WithErrorDetails(err.Error())
        return
    }

    resp.SetResult(http.StatusOK, models)
}
```
