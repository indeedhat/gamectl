import { fileSize, trafficSpeed, duration, once, loadCss, percentage } from '/assets/js/util.js';

class ResourceMonitor
{
    constructor(selector)
    {
        this.$monitor = document.querySelector(selector);
        this.$data = this.$monitor.querySelector('.data');

        this._stream = null;

        this.$data.innerHTML = loadingTemplate();

        this._importCss();
        this._setupEventSource();
        this._destructor();
    }

    _destructor()
    {
        window.onbeforeunload = () => {
            this._stream.close();
        };
    }

    _importCss()
    {
        once(
            "resource-monitor", 
            () => loadCss('/assets/css/resource-monitor.css')
        );
    }

    _setupEventSource()
    {
        this._stream = new EventSource("/api/performance");

        this._stream.onerror = this._handleError.bind(this);
        this._stream.addEventListener("message", this._handleMessage.bind(this));
    }


    _handleError(e)
    {
        console.error("Resource Monitor:", e);
    }

    _handleMessage({ data })
    {
        this.$data.innerHTML = buildTemplate(JSON.parse(data));
    }

}

const loadingTemplate = () => {
    return `
        <table>
            <tr>
                <td>Loading System Resources...</td>
            </tr>
        </table>
    `;
}

const buildTemplate = ({ uptime, memory, cpu, network, mount }) => {
    let coreTemplate = "";
    let networkTemplate = "";
    let mountTemplate = "";
    let cpuTotal = 0;
    let cpuIdle = 0;

    for (let [ key, core ] of Object.entries(sortCores(cpu))) {
        cpuTotal += core.total;
        cpuIdle += core.idle;

        coreTemplate += `
            <div class="core">
                <strong>${key}:</strong>
                ${percentage(core.total, core.total - core.idle)}
            </div>`;
    }

    for (let [ key, intf ] of Object.entries(network)) {
        networkTemplate += `
            <div class="intf">
                <strong>${key}:</strong>
                <div>Tx: ${trafficSpeed(intf.tx)}</div>
                <div>Rx: ${trafficSpeed(intf.rx)}</div>
            </div>
        `;
    }

    for (let [ key, mnt ] of Object.entries(mount)) {
        mountTemplate += `
            <div class="mnt">
                <strong>${key}:</strong>
                ${fileSize(mnt.used)} / ${fileSize(mnt.total)}
            </div>
        `;
    }

    return `
        <table>
            <tr>
                <th>Uptime</th>
                <td>${duration(uptime)}</td>
            </tr>
            <tr>
                <th>Memory</th>
                <td>${fileSize(memory.used)} / ${fileSize(memory.total)}</td>
            </tr>
            <tr>
                <th>CPU</th>
                <td>
                    <div>
                        <strong>Total:</strong>
                        ${percentage(cpuTotal, cpuTotal - cpuIdle)}
                    </div>
                    ${coreTemplate}
                </td>
            </tr>
            <tr>
                <th>Network</th>
                <td>${networkTemplate}</td>
            </tr>
            <tr>
                <th>Mounts</th>
                <td>${mountTemplate}</td>
            </tr>
        </table>
    `;
};

const sortCores = cpu => {
    return Object.keys(cpu)
        .sort((a, b) => {
            if (a.length > b.length) {
                return 1;
            }

            return a > b ? 1 : 0;
        })
        .reduce((object, key) => {
            object[key] = cpu[key];
            return object;
        }, {})
};


export default ResourceMonitor;