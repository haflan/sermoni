
/**
 * Create URL parameters from a JSON object
 */
function params(parameterObject) {
    const parameters = new URLSearchParams();
    for (let key in parameterObject) {
        let value = parameterObject[key];
        if (Array.isArray(value)) {
            for (let i = 0; i < value.length; i++) {
                parameters.append(key, value[i]);
            }
        } else {
            parameters.set(key, value);
        }
    }
    return parameters.toString();
}

var csrfToken;

// Send a request to server. Takes request description, jsonRequest,
// on the following format
// {
//      url: "/services" ,
//      method: "POST",
//      data: { "newId": 8234 }
//      expectedStatus: 201,
//      error: errData => { console.log(errData); },
//      success: successData => { console.log(successData); }
// }
function request(jsonRequest) {
    if (!jsonRequest.url) {
        console.error("No URL provided in the request object");
        return;
    }
    let xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (!jsonRequest.expectedStatus) {
            jsonRequest.expectedStatus = 200
        }
        if (this.readyState === XMLHttpRequest.DONE) {
            const data = this.responseText ? JSON.parse(this.responseText) : {};
            if (this.status !== jsonRequest.expectedStatus) {
                if (jsonRequest.error) {
                    jsonRequest.error(this.status, data);
                } else {
                    console.error(this.status + ": " + this.responseText);
                }
            } else if (jsonRequest.success) {
                jsonRequest.success(data);
            }
        }
    };
    if (jsonRequest.method) {
        xhttp.open(jsonRequest.method, jsonRequest.url, true);
    } else {
        xhttp.open("GET", jsonRequest.url, true);
    }
    xhttp.setRequestHeader("X-Csrf-Token", csrfToken);
    if (jsonRequest.data) {
        xhttp.setRequestHeader("Content-Type", "application/json");
        xhttp.send(JSON.stringify(jsonRequest.data));
    } else {
        xhttp.send();
    }
}

function init(successHandler, errorHandler) {
    request({
        url: "/init",
        success: (data) => {
            csrfToken = data.csrftoken;
            if (successHandler) {
                successHandler(data);
            }
        },
        error: errorHandler
    });
}

function login(passphrase, successHandler, errorHandler) {
    request({
        url: "/login",
        method: "POST",
        data: { passphrase: passphrase },
        success: successHandler,
        error: errorHandler
    });
}

export default {
    init,
    login,
}
