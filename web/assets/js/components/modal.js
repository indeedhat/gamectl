import { html, loadCss, once } from '/assets/js/util.js';

class Modal
{
    constructor(config = DefaultConfig)
    {
        this._importCss();

        this._initAlpine(this._onInit.bind(this, config));
    }

    _onInit(config, modal)
    {
        this.modal = modal;
        this.modal.config = {
            ...this.modal.config,
            ...config
        };

        this.modal.open();
    }

    _initAlpine(then)
    {
        if ("undefined" !== typeof window.modal) {
            then(window.modal);
            return 
        }

        document.body.append(...html`
            <div id="modal-wrapper" x-data="modal" x-show="visible" style="display:block">
                <div id="modal-cover"></div>
                <div id="modal" @click.outside="close">
                    <div id="modal-title">&gt;<span x-text="config.title"></span></div>
                    <div id="modal-body" x-html="config.content"></div>
                </div>
            </div>
        `);

        window.addEventListener("alpine:init", function() {
            Alpine.data("modal", () => ({
                config: DefaultConfig,
                visible: false,

                init() {
                    window.modal = this;
                    then(this);
                },

                open() {
                    this.visible = true;

                    if ("function" === typeof this.config.onOpen) {
                        this.config.onOpen(this);
                    }
                },

                close() {
                    this.visible = false;

                    if ("function" === typeof this.config.onClose) {
                        this.config.onClose(this);
                    }
                }
            }));
        });
    }

    _importCss()
    {
        once("moduleCss", () => {
            loadCss("/assets/css/modal.css");
        });
    }
}

const DefaultConfig = {
    title: "Modal",
    content: null,
    onOpen: (modal) => console.log("modal opened", modal),
    onClose: (modal) => console.log("modal closed", modal),
};

export default Modal;
export {
    Modal,
    DefaultConfig
}
