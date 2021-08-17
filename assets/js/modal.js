import { once, html } from "/assets/js/util.js";

class Modal
{
    constructor(config = DefaultConfig) 
    {
        this.config = config;

        this._importCss();
        this._createModalElements();
        this._watchElements();
        this._setupBinds();
    }

    /** 
     * Open the modal
     *
     * When opened the modal content will be refreshed from the template given (element or text)
     *
     * if set the onOpen closure will be called
     *
     * @return void
     */
    open()
    {
        this.$title.innerText = this.config.title;
        this.$body.innerHTML = "";

        if (this.$contentTemplate) {
            this.$body.innerHTML = this.$contentTemplate.innerHTML;
        } else if ("string" === typeof this.config.content) {
            this.$body.append(...html(this.config.content));
        } else if ("function" === typeof this.config.conent[Symbol.iterator]) {
            this.$body.append(...this.config.conent);
        } else if (this.config.content) {
            this.$body.append(this.config.content);
        }

        this.$wrapper.style.display = "block";

        if ("function" === typeof this.config.onOpen) {
            this.config.onOpen();
        }
    }

    /**
     * Close the modal
     *
     * If set the onClose closure will be called
     *
     * @return void
     */
    close()
    {
        this.$wrapper.style.display = "none";

        if ("function" === typeof this.config.onClose) {
            this.config.onClose();
        }
    }

    _importCss()
    {
        once("moduleCss", () => {
            document.head.append(...html`<link type="text/css" rel="stylesheet" href="/assets/css/modal.css" />`);
        });
    }

    _createModalElements()
    {
        once("moduleHtml", () => {
            document.body.append(...html`
                <div id="modal-wrapper">
                    <div id="modal-cover"></div>
                    <div id="modal">
                        <div id="modal-title">&gt;<span></span></div>
                        <div id="modal-body"></div>
                    </div>
                </div>
            `);
        });

    }

    _setupBinds()
    {
        this.$cover.onclick = this.close.bind(this);
    }

    _watchElements()
    {
        this.$wrapper = document.getElementById("modal-wrapper");
        this.$cover = document.getElementById("modal-cover");
        this.$modal = document.getElementById("modal");
        this.$title = document.querySelector("#modal-title span");
        this.$body = document.getElementById("modal-body");

        if (this.config.contentSelector) {
            this.$contentTemplate = document.querySelector(this.config.contentSelector);
        }
    }
}


const DefaultConfig = {
    title: "Modal",
    content: null,
    contentSelector: "#modal-content",
    onOpen: () => console.log("modal opened"),
    onClose: () => console.log("modal closed"),
};

export default Modal;

export {
    Modal,
    DefaultConfig
};
