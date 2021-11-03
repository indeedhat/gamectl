import { fileSize, trafficSpeed, duration, once, loadCss, percentage } from '/assets/js/util.js';

/**
 * Setup the Resource monitor components
 *
 * @async
 */
const createResourceMonitor = async () => {
    once("resource-monitor", () => loadCss('/assets/css/resource-monitor.css'));

    window.addEventListener("alpine:init", function() {
        Alpine.data("resourceMonitor", () => ({
            loading: true,
            uptime: "string",
            memory: {},
            cpu: {},
            mounts: {},
            network: {},

            init() {
                new StreamHandler(this);
            }
        }));
    });
};

class StreamHandler
{
    constructor(monitorData)
    {
        this.monitor = monitorData;
        this.stream  = null;

        this.setupEventSource();

        window.onbeforeunload = () => {
            this.stream.close();
        };
    }

    setupEventSource()
    {
        this.stream = new EventSource("/api/performance");

        this.stream.onerror = e => console.error("Resource Monitor:", e, this.monitor);
        this.stream.addEventListener("message", this.handleMessage.bind(this));
    }

    handleMessage({ data })
    {
        let { uptime, memory, cpu, network, mount } = JSON.parse(data);

        this.monitor.loading = false;
        this.monitor.uptime  = duration(uptime);
        this.monitor.cpu     = buildCpu(cpu);
        this.monitor.memory  = buildMemory(memory);
        this.monitor.mounts  = buildMounts(Object.entries(mount));
        this.monitor.network = buildNetwork(Object.entries(network));
    }
}

const buildCpu = cpu => {
    let cpuTotal = 0;
    let cpuIdle = 0;
    let cores = [];

    for (let [ key, core ] of Object.entries(sortCores(cpu))) {
        cpuTotal += core.total;
        cpuIdle += core.idle;

        let corePercent = percentage(core.total, core.total - core.idle);

        cores.push({
            key,
            percent: corePercent,
            color: `rgb(${cpuTemp(corePercent)})`
        });
    }

    let cpuPercent = percentage(cpuTotal, cpuTotal - cpuIdle);
    return {
        percent: cpuPercent,
        color: `rgb(${cpuTemp(cpuPercent)})`,
        cores
    };
};

const buildMemory = memory => {
    let memoryPercent = percentage(memory.total, memory.used);

    return {
        text: fileSize(memory.used) + " / " + fileSize(memory.total),
        width: memoryPercent,
        color: `rgb(${cpuTemp(memoryPercent)})`
    };
}

const buildMounts = mount => mount.map(([ key, mnt ]) => {
    let mountPercent = percentage(mnt.total, mnt.used);

    return {
        key,
        text: `${fileSize(mnt.used)} / ${fileSize(mnt.total)}`,
        width: mountPercent,
        color: `rgb(${cpuTemp(mountPercent)})`
    };
});

const buildNetwork = network => network.map(([ key, intf ]) => ({
    key,
    tx: trafficSpeed(intf.tx),
    rx: trafficSpeed(intf.rx)
}));


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

export default createResourceMonitor;
