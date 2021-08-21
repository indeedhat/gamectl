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

/**
 * Load a css file by adding the link element to the head
 *
 * @param string src
 *
 * @return void
 */
const loadCss = href => {
    console.log("loading css:", href);
    document.head.append(
        ...html(`<link type="text/css" rel="stylesheet" href="${href}" />`)
    );
};

/**
 * Load a css file by adding the link element to the head
 *
 * @param string src
 *
 * @return void
 */
const loadJs = src => {
    console.log("loading js:", src);

    const script = document.createElement("script");
    script.src = src;
    document.head.append(script);
};

export {
    once,
    html,
    objectToForm,
    loadCss,
    loadJs
};
