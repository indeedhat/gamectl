import { once, loadJs } from '/assets/js/util.js';

once(
    "js-downloader", 
    () => loadJs("/assets/js/lib/js-file-downloader.min.js")
);

const initDownloadButton = (key, $elem) => {
    $elem.onclick = _buildDownloadEventHandler(key, $elem);
};

const _buildDownloadEventHandler = (key, $elem) => {
    return async function() {
        if (!confirm("The game server will be shut down while the backup us processed!\nDo you wish to continue?")) {
            return;
        }

        let download = new jsFileDownloader({
            url: `/api/apps/${key}/download`,
            autoStart: false
        });

        $elem.innerHTML = "Please Wait...";

        try {
            await download.start();
            $elem.innerHTML = "Download World";
        } catch (e) {
            $elem.innerHTML = "Download Failed!";
        }
    };
};

export {
    initDownloadButton
};
