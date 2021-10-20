import { objectToForm } from "/assets/js/util.js";

/**
 * Perform a GET request with the expectiation of a JSON response
 *
 * @param string url
 * @param object data
 *
 * @return Promise
 */
const get = (url, data = {}) => {
    if (Object.keys(data).length) {
        url = `${url}?${new URLSearchParams(data).toString()}`;
    }

    let status = 0;

    return fetch(url)
        .then(r => {
            status = r.status;
            return r.json();
        })
        .then(json => [status, json]);
};

/**
 * Perform a POST request with the expectiation of a JSON response
 *
 * @param string url
 * @param object data
 *
 * @return Promise
 */
const post = (url, data = {}) => {
    let status = 0;

    return fetch(url, { method: "POST", body: objectToForm(data) })
        .then(r => {
            status = r.status;
            return r.json();
        })
        .then(json => [status, json]);
};

export {
    get,
    post
}
