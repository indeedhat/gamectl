/**
 * Run a piece of js code once
 *
 * after it has been run it will be tracked in the window object, should the code be called again
 * then it will be skipped
 * 
 * @param string key
 * @param function fn
 *
 * @return void
 */
const once = (key, fn) => {
    if ("undefined" === typeof window.onceRegister) {
        window.onceRegister = [];
    }

    if (~window.onceRegister.indexOf(key)) {
        return;
    }

    fn();
    window.onceRegister.push(key);
};

/**
 * convert a string to html elements
 *
 * @param string htmlString
 * 
 * @return NodeList
 */
const html = htmlString => {
    let template = document.createElement("template");

    template.innerHTML = htmlString;

    return template.content.childNodes;
};

/**
 * Convert an object for a form instance
 *
 * @param object data
 *
 * @return FormData
 */
const objectToForm = data => {
    let form = new FormData();

    for (let key in data) {
        form.append(key, data[key]);
    }

    return form;
};

export {
    once,
    html,
    objectToForm
};