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
    const script = document.createElement("script");
    script.src = src;
    document.head.append(script);
};

const duration = unix => {
    const days = Math.floor(unix / 86400);
    const hours = Math.floor((unix % 86400) / 3600);
    const minutes = Math.floor((unix % 3600) / 60);
    const seconds = Math.floor(unix % 60);

    return `${days} days, ${hours} hours, ${minutes} minuts, ${seconds} seconds`;
};

const fileSize = bytes => {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const dm = 2;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
};

const trafficSpeed = bytes => {
    if (bytes === 0) return '0 bps';

    const k = 1024;
    const dm = 2;
    const sizes = ['bps', 'kbs', 'mbs', 'gbs'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
};

export {
    once,
    html,
    objectToForm,
    loadCss,
    loadJs,
    duration,
    fileSize,
    trafficSpeed
};
