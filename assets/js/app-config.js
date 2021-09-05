import { once, html, loadCss, loadJs, escape } from "/assets/js/util.js";
import { get, post } from "/assets/js/request.js";
import { Modal } from "/assets/js/modal.js";

class AppConfigController 
{
    constructor(appName, appKey)
    {
        this._importCss();

        this.appName = appName;
        this.appKey  = appKey;

        this.configFiles = [];
        this.logFiles    = [];

        this.logSource = null;

        this.$app = document.querySelector(`.serverStatus[data-key=${appKey}]`);
        if (!this.$app) {
            return;
        }

        this.$configButton = this.$app.querySelector(".controller.config");
        this.$logButton = this.$app.querySelector(".controller.logs");

        this._setupConfig();
        this._setupLogs()
        this._loadCodeMirrorModes();
    }

    _loadCodeMirrorModes()
    {
        for (let key in this.configFiles) {
            let mode = this.configFiles[key].mode;
            if (!mode) {
                continue;
            }

            if (mode === "json") {
                mode = "javascript";
            }

            once(`CodeMirror.${mode}`, () => {
                loadJs(`/assets/js/lib/codemirror/mode/${mode}/${mode}.js`);
            });
        }
    }

    _setupConfig()
    {
        try {
            this.configFiles = JSON.parse(this.$app.dataset.config);
            if (!Object.keys(this.configFiles).length) {
                console.log(this.configFiles, this.configFiles.length)
                throw new Error();
            }

            this.$configButton.removeAttribute("disabled");
            this._initializeConfigModal();
        } catch (e) {
            this.$configButton.setAttribute("disabled", true);
        }
    }

    _setupLogs()
    {
        try {
            this.logFiles = JSON.parse(this.$app.dataset.logs);
            if (!Object.keys(this.logFiles).length) {
                throw new Error();
            }

            this.$logButton.removeAttribute("disabled");
            this._initialzeLogsModal();
        } catch (e) {
            this.$logButton.setAttribute("disabled", true);
        }
    }

    _importCss()
    {
        once("appConfigCss", () => {
            loadCss("/assets/css/app-config.css");
        });
    }

    _initializeConfigModal()
    {
        this.configModal = new Modal({
            title: `${this.appName} Config`,
            content: configModalTemplate(this.appKey, map(this.configFiles, info => info.description)),
            onOpen: this._handleConfigModalOpen.bind(this)
        });

        this.$configButton.onclick = () => {
            this.configModal.open();
        };

    }

    _initialzeLogsModal()
    {
        this.logsModal = new Modal({
            title: `${this.appName} Logs`,
            content: logsModalTemplate(this.appKey, this.logFiles),
            onOpen: this._handleLogsModalOpen.bind(this),
            onClose: () => {
                if (this.logSource) {
                    this.logSource.close();
                    this.logSource = null;
                }
            }
        });

        this.$logButton.onclick = () => {
            this.logsModal.open();
        };
    }

    _handleConfigModalOpen()
    {
        let $back = this.configModal.$body.querySelector("a.back");
        let $configList = this.configModal.$body.querySelector("div.config-list");
        let $listEntries = $configList.querySelectorAll("div.block-entry");

        let $configLoading = this.configModal.$body.querySelector("div.config-loading");

        let $configForm = this.configModal.$body.querySelector("div.config-form");
        let $formError = $configForm.querySelector(".error");
        let $formAlert = $configForm.querySelector(".alert");
        let $textArea = $configForm.querySelector("textarea");

        $back.onclick = () => {
            $configForm.style.display    = "none";
            $configLoading.style.display = "none";
            $configList.style.display    = "block";
        };

        $listEntries.forEach(element => element.onclick = async e => {
            const target = e.target.className === "block-entry" ? e.target : e.target.parentNode;
            const configKey = target.dataset.key;

            $formError.innerHTML = "";
            $formAlert.innerHTML = "";

            let $existingCodeMirror = $configForm.querySelector(".CodeMirror")
            let codeMirrorInstance = null;
            if ($existingCodeMirror) {
                $existingCodeMirror.remove();
            }

            $configForm.style.display    = "none";
            $configLoading.style.display = "block";
            $configList.style.display    = "none";

            $configForm.querySelector("form").onsubmit = async e => {
                e.preventDefault();

                $formAlert.innerHTML = "";
                $formError.innerHTML = "";
                codeMirrorInstance.save();

                try {
                    let [status, response] = await post(`/api/apps/${this.appKey}/config/${configKey}`, {data: $textArea.value});

                    if (status === 404) {
                        $formError.innerHtml = "Config not found";
                    } else if (status == 500) {
                        $formError.innerHtml = "Update Failed";
                    } else if (!response.outcome) {
                        $formError.innerHtml = response.message || "Unknown Error";
                    } else {
                        $formAlert.innerHTML = "Config Updated";
                    }
                } catch (e) {
                    $formError.innerHTML = "Unknwon Error";
                    console.error(e);
                }
            };

            try {
                let [_, response] = await get(`/api/apps/${this.appKey}/config/${configKey}`);

                if (!response.outcome) {
                    throw new Error("failed");
                }

                $configForm.querySelector("h2").innerHTML = configKey;

                $textArea.value = response.file;
                codeMirrorInstance = this._initCodeMirror($textArea, this.configFiles[configKey].mode);

                $configLoading.style.display = "none";
                $configForm.style.display    = "block";
            } catch (e) {
                $configLoading.style.display = "none";
                $configList.style.display    = "block";
                console.error(e);
            }
        });
    }

    _initCodeMirror(element, mode)
    {
        if (mode === "json") {
            mode = {
                name: "javascript",
                json: true
            };
        }

        let cm = CodeMirror.fromTextArea(element, {
            lineNumbers: true,
            autoRefresh: true,
            lineWrapping: true,
            theme: "darcula",
            mode
        });

        setTimeout(() => cm.refresh());

        return cm;
    }

    _handleLogsModalOpen()
    {
        let $logLoading = this.logsModal.$body.querySelector("div.log-loading");

        let $logWrapper = this.logsModal.$body.querySelector("div.log-wrapper");
        let $logs = $logWrapper.querySelector("pre.logs");

        let $back = this.logsModal.$body.querySelector("a.back");
        let $logsList = this.logsModal.$body.querySelector("div.log-list");
        let $listEntries = $logsList.querySelectorAll("div.block-entry");

        $back.onclick = () => {
            $logWrapper.style.display = "none";
            $logLoading.style.display = "none";
            $logsList.style.display   = "block";

            if (this.logSource) {
                this.logSource.close();
                this.logSource = null;
            }
        };

        $listEntries.forEach(element => element.onclick = async e => {
            const target = e.target.className === "block-entry" ? e.target : e.target.parentNode;
            const logKey = target.dataset.key;

            $logs.innerHTML = "";

            $logWrapper.style.display = "none";
            $logLoading.style.display = "block";
            $logsList.style.display   = "none";

            try {
                this.logSource = new EventSource(`/api/apps/${this.appKey}/logs/${logKey}`);
                this.logSource.addEventListener("message", e => {
                    scrollToBottom($logs, () => {
                        $logs.innerHTML += escape(e.data);
                    });
                }, false);

                this.logSource.addEventListener("keep-alive", e => {
                    console.log(`keep-alive: ${logKey}`);
                }, false);

                this.logSource.onopen = () => {
                    $logWrapper.querySelector("h2").innerHTML = logKey;

                    $logLoading.style.display = "none";
                    $logsList.style.display   = "none";
                    $logWrapper.style.display = "block";
                    $logs.innerHTML = "";
                };
                
                this.logSource.onerror = e => {
                    $logWrapper.style.display = "none";
                    $logLoading.style.display = "none";
                    $logsList.style.display   = "block";
                    console.error(e);
                };
            } catch (e) {
                $logLoading.style.display = "none";
                $logsList.style.display   = "block";
                console.error(e);
            }
        });
    }
}

/**
 * If the element is currently scrolled to the bottom it will reestablish that scroll position
 * once the given closure is called
 *
 * @param DomElement element
 * @param function fn
 *
 * @return void
 */
const scrollToBottom = (element, fn) => {
    const scrollLeeway = 10;
    let startingPosition = element.offsetHeight + element.scrollTop;
    let startingScrollHeight = element.scrollHeight;

    fn();

    if (startingPosition + scrollLeeway >= startingScrollHeight) {
        element.scrollTop = element.scrollHeight;
    }
};

/**
 * Build up a template for the config modal based on the app config
 *
 * @param string appKey
 * @param array configFiles
 *
 * @return string
 */
const configModalTemplate = (appKey, configFiles) =>  `
    <div class="app-config">
        <div class="config-loading">Loading...</div>

        <div class="config-list">
            <h2>Server Configuration Files</h2>
            <div class="error">
                If you dont know what your doing please leave well enough alone
                <br />
                <br />
                there is no validation on updating config and uploading an invalid file may
                break the server.
            </div>
            <div class="block-list">
                ${buildBlocks(configFiles)}
            </div>
        </div>

        <div class="config-form">
            <a href="javascript: void(0);" class="back">&lt; Back</a>
            <h2></h2>
            <form method="post" class="full-width">
                <div class="error"></div>
                <div class="alert"></div>

                <div class="field">
                    <label>Config Body</label>
                    <textarea name="data" rows="32"></textarea>
                </div>

                <div class="group">
                    <input type="submit" value="Update Config" />
                </div>
            </form>
        </div>
    </div>
`;

/**
 * Build up a template for the config modal based on the app config
 *
 * @param string appKey
 * @param array logFiles
 *
 * @return string
 */
const logsModalTemplate = (appKey, logFiles) =>  `
    <div class="app-logs">
        <div class="log-loading">Loading...</div>

        <div class="log-list">
            <h2>Server Logs</h2>
            <div class="block-list">
                ${buildBlocks(logFiles)}
            </div>
        </div>

        <div class="log-wrapper">
            <a href="javascript: void(0);" class="back">&lt; Back</a>
            <h2></h2>
            <pre class="logs"></pre>
        </div>
    </div>
`;

/**
 * Build the blocks for the lists
 *
 * @param array entries
 * 
 * @return string
 */
const buildBlocks = entries => {
    let blocks = "";
    for (let key in entries) {
        blocks += `
            <div class="block-entry" data-key="${key}">
                ${key}
                <div class="description">${entries[key]}</div>
            </div>`;
    }

    return blocks;
};

/**
 * Map an object
 *
 * @param object object
 * @param function fn
 *
 * @return object
 */
const map = (object, fn) => {
    let out = {};
    for (let key in object) {
        out[key] = fn(object[key]); 
    }

    return out;
};

export default AppConfigController;
