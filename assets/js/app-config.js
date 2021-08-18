import { once, html } from "/assets/js/util.js";
import { get, post } from "/assets/js/request.js";
import { Modal } from "/assets/js/modal.js";

class AppConfigController 
{
    constructor(appName, appKey)
    {
        this._importCss();

        this.appName = appName;
        this.appKey  = appKey;

        this.$app = document.querySelector(`.serverStatus[data-key=${appKey}]`);
        if (!this.$app) {
            return;
        }

        this.$configButton = this.$app.querySelector(".controller.config");
        if (!this._loadConfig()) {
            this.$configButton.classList.add("disabled");
            return;
        }

        this._setupModal();
    }

    _loadConfig()
    {
        try {
            this.configKeys = JSON.parse(this.$app.dataset.config);
        } catch (e) {
            return false;
        }

        return true;
    }

    _importCss()
    {
        once("appConfigCss", () => {
            document.head.append(
                ...html`<link type="text/css" rel="stylesheet" href="/assets/css/app-config.css" />`
            );
        });
    }

    _setupModal()
    {
        this.modal = new Modal({
            title: `${this.appName} Config`,
            content: modalTemplate(this.appKey, this.configKeys),
            onOpen: this._handleModalOpen.bind(this)
        });

        this.$configButton.onclick = () => {
            this.modal.open();
        };
    }

    _handleModalOpen()
    {
        let $back = this.modal.$body.querySelector("a.back");
        let $configList = this.modal.$body.querySelector("div#config-list");
        let $listEntries = $configList.querySelector("div.block-entry");

        let $configLoading = this.modal.$body.querySelector("div#config-loading");

        let $configForm = this.modal.$body.querySelector("div#config-form");
        let $formError = $configForm.querySelector(".error");
        let $formAlert = $configForm.querySelector(".alert");
        let $textArea = $configForm.querySelector("textarea");

        $back.onclick = () => {
            $configForm.style.display    = "none";
            $configLoading.style.display = "none";
            $configList.style.display    = "block";
        };

        $listEntries.onclick = async e => {
            const configKey = e.target.innerText.trim();
            $formError.innerHTML = "";
            $formAlert.innerHTML = "";

            $configForm.style.display    = "none";
            $configLoading.style.display = "block";
            $configList.style.display    = "none";

            $configForm.querySelector("form").onsubmit = async e => {
                e.preventDefault();

                $formAlert.innerHTML = "";
                $formError.innerHTML = "";

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

                $textArea.value = response.file;
                $configLoading.style.display = "none";
                $configForm.style.display    = "block";
            } catch (e) {
                $configLoading.style.display = "none";
                $configList.style.display    = "block";
                console.error(e);
            }
        };
    }
}

/**
 * Build up a template for the modal based on the app config
 *
 * @param string appKey
 * @param array configKeys
 *
 * @return string
 */
const modalTemplate = (appKey, configKeys) => {
    let blocks = "";
    for (let i in configKeys) {
        blocks += `<div class="block-entry">${configKeys[i]}</div>`;
    }

    return `
        <div id="app-config">
            <div id="config-loading">Loading...</div>

            <div id="config-list">
                <h2>Server Configuration Files</h2>
                <div class="error">
                    If you dont know what your doing please leave well enough alone
                    <br />
                    <br />
                    there is no validation on updating config and uploading an invalid file may
                    break the server.
                </div>
                <div class="block-list">
                    ${blocks}
                </div>
            </div>

            <div id="config-form">
                <a href="javascript: void(0);" class="back">&lt; Back</a>
                <form method="post">
                    <div class="error"></div>
                    <div class="alert"></div>

                    <div class="field">
                        <label>Config Body</label>
                        <textarea name="data" rows="16"></textarea>
                    </div>

                    <div class="group">
                        <input type="submit" value="Update Config" />
                    </div>
                </form>
            </div>
        </div>
    `;
};

export default AppConfigController;
