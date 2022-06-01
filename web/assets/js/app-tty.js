import { once, loadCss } from "/assets/js/util.js";
import { Modal } from "/assets/js/modal.js";

class AppTTY 
{
    constructor(appName, appKey)
    {
        this._importCss();

        this.socket = null;
        this.appName = appName;
        this.appKey  = appKey;

        this.$app = document.querySelector(`.server[data-key=${appKey}]`);
        if (!this.$app) {
            return;
        }

        this.$ttyButton = this.$app.querySelector(".config");

        this._setupTty();
    }

    _setupLogs()
    {
        try {
            this.$ttyButton.removeAttribute("disabled");
            this._initializeModal();
        } catch (e) {
            this.$logButton.setAttribute("disabled", true);
        }
    }

    _importCss()
    {
        once("appConfigCss", () => loadCss("/assets/css/app-config.css"));
    }

    _initialzeLogsModal()
    {
        this.modal = new Modal({
            title: `${this.appName} TTY`,
            content: modalTemplate(this.appKey),
            onOpen: this._handleModalOpen.bind(this),
            onClose: () => {
                if (this.socket) {
                    this.socket.close();
                    this.socket = null;
                }
            }
        });

        this.$ttyButton.onclick = () => {
            this.modal.open();
        };
    }

    _handleModalOpen()
    {
        let $loading = this.modal.$body.querySelector("div.loading");
        let $wrapper = this.modal.$body.querySelector("div.wrapper");
        let $shell = $logWrapper.querySelector("pre.shell");
        let $input = this.modal.$body.querySelector("input.tty-input");

        $shell.innerHTML = "";
        $wrapper.style.display = "none";
        $loading.style.display = "block";

        const writeLog = txt => {
            const doScroll = $shell.scrollTop > $shell.scrollHeight - $shell.clientHeight - 1;
            let item = document.createElement("div");
            item.innerHTML = txt;

            $shell.appendChild(item);
            if (doScroll) {
                $shell.scrollTop = $shell.scrollHeight - $shell.clientHeight;
            }
        };

        this.socket = new WebSocket(`ws://${document.location.host}/api/apps/${this.appKey}/tty`);
        this.socket.onclose = () => {
            writeLog(`<strong>Connection Closed</strong>`);
        };
        this.socket.onmessage = e => {
            let messages = e.data.split('\n');
            for (let i = 0; i < messages.length; i++) {
                writeLog(ansispan(messages[i]));
            }
        };

        $input.onkeydown = e => {
            // eat tabs
            if (e.keyCode == 9) {
                e.preventDefault();
                return false;
            }
        };
        $input.onkeyup = e => {
            e.preventDefault();

            if (!this.socket) {
                return false;
            }

            // enter
            let code = e.keyCode;
            if (e.keyCode != 13) {
                $input.value = "";
                code = 10
            } else if (code >= 65 && code <= 90 && !e.shiftKey) {
                code += 32
            }

            this.socket.send(String.fromCharCode(code));

            return false;
        }
    }
}


/**
 * Build up a template for the config modal based on the app config
 *
 * @param string appKey
 * @param array logFiles
 *
 * @return string
 */
const modalTemplate = () =>  `
    <div class="app-tty">
        <div class="loading">Loading...</div>
        <div class="wrapper">
            <pre class="shell"></pre>
            <input class="tty-input" type="text" />
        </div>
    </div>
`;


const ansispan = str => {
    Object.keys(colors).forEach(function (ansi) {
        let span = '<span style="color: ' + colors[ansi] + '">';

        str = str.replace(new RegExp(`\033\\[${ansi}m`, 'g'), span)
            .replace(new RegExp(`\033\\[0;${ansi}m`, 'g'), span);
    });

    return str.replace(/\033\[1m/g, '<b>').replace(/\033\[22m/g, '</b>')
        .replace(/\033\[3m/g, '<i>').replace(/\033\[23m/g, '</i>')
        .replace(/\033\[m/g, '</span>')
        .replace(/\033\[0m/g, '</span>')
        .replace(/\033\[39m/g, '</span>');
};

const colors = {
    '30': 'black',
    '31': 'red',
    '32': 'green',
    '33': 'yellow',
    '34': 'blue',
    '35': 'purple',
    '36': 'cyan',
    '37': 'white'
};

export default AppTTY
