import { post, get } from '/assets/js/request.js';
import { minDuration, once, loadJs } from '/assets/js/util.js';

class AppControls
{
    constructor(appKey)
    {
        this.appKey  = appKey;

        this.$app = document.querySelector(`.serverStatus[data-key=${appKey}]`);
        if (!this.$app) {
            return;
        }

        this.$startStopButton = this.$app.querySelector(".controller.startStop");
        this.$restartButton = this.$app.querySelector(".controller.restart");
        this.$downloadButton = this.$app.querySelector(".controller.download");

        this.$status = this.$app.querySelector(".dataTile.status");
        this.$players = this.$app.querySelector(".dataTile.players .current");
        this.$uptime = this.$app.querySelector(".dataTile.uptime");

        this.updateStatus();
        setTimeout(this.updateStatus.bind(this), 30000);

        this._includeJsDownloader();
        this._initDownloadButton();
    }

    async updateStatus()
    {
        this._handleUpdateStatus()
    }

    _includeJsDownloader()
    {
        once(
            "js-downloader", 
            () => loadJs("/assets/js/lib/js-file-downloader.min.js")
        );
    }

    _initDownloadButton()
    {
        if (!this.$downloadButton.disabled) {
            this.$downloadButton.onclick = this._handleDownload.bind(this);
        }
    }

    async _handleStart()
    {
        this.$startStopButton.innerHTML = "Starting...";
        this.$startStopButton.setAttribute("dsiabled", true);

        try {
            let [, json] = await post(`api/apps/${this.appKey}/start`, {});
            if (!json.outcome) {
                throw new Error();
            }

            await this.updateStatus();
        } catch (e) {
            this.$startStopButton.innerHTML ="Start Failed!";
        }

        this.$startStopButton.removeAttribute("dsiabled");
    }

    async _handleStop()
    {
        this.$startStopButton.innerHTML = "Stopping...";
        this.$startStopButton.setAttribute("dsiabled", true);

        try {
            let [, json] = await post(`api/apps/${this.appKey}/stop`, {});
            if (!json.outcome) {
                throw new Error();
            }

            this.updateStatus();
        } catch (e) {
            this.$startStopButton.innerHTML ="Stop Failed!";
        }

        this.$startStopButton.removeAttribute("dsiabled");
    }

    async _handleRestart()
    {
        this.$restartButton.innerHTML = "Restarting...";
        this.$restartButton.setAttribute("dsiabled", true);

        try {
            let [, json] = await post(`api/apps/${this.appKey}/restart`, {});
            if (!json.outcome) {
                throw new Error();
            }

            this.updateStatus();
        } catch (e) {
            this.$restartButton.innerHTML ="Restart Failed!";
        }

        this.$restartButton.removeAttribute("dsiabled");
    }

    async _handleUpdateStatus()
    {
        this.$status.classList.remove("offline", "online");
        this.$restartButton.onclick = null;
        this.$startStopButton.onclick = null;
        this.$restartButton.setAttribute("disabled", true);

        try {
            let [, json] = await get(`/api/apps/${this.appKey}`);
            if (!json.outcome) {
                throw new Error();
            }

            this.$players.innerHTML = -1 === ~~json.status.users ? "N/A" : json.status.users;
            this.$uptime.innerHTML = minDuration(json.status.uptime);

            if (!json.status.online) {
                this.$status.innerHTML = "Offline";
                this.$status.classList.add("offline");

                this.$startStopButton.innerHTML = "Start";
                this.$startStopButton.onclick = this._handleStart.bind(this);
                return;
            }

            this.$status.innerHTML = "Online";
            this.$status.classList.add("online");
            this.$startStopButton.innerHTML = "Stop";
            this.$startStopButton.onclick = this._handleStop.bind(this);
            this.$restartButton.onclick = this._handleRestart.bind(this);
            this.$restartButton.removeAttribute("disabled");
        } catch (e) {
            console.log(e)
            this.$status.innerHTML = "UNAVAILABLE";
            this.$status.classList.add("offline");
        }
    }

    async _handleDownload()
    {
        if (!confirm("The game server will be shut down while the backup us processed!\nDo you wish to continue?")) {
            return;
        }

        let download = new jsFileDownloader({
            url: `/api/apps/${this.appKey}/download`,
            autoStart: false
        });

        this.$downloadButton.innerHTML = "Please Wait...";

        try {
            await download.start();
            this.$downloadButton.innerHTML = "Download World";
        } catch (e) {
            this.$downloadButton.innerHTML = "Download Failed!";
        }
    }
}

export default AppControls;
