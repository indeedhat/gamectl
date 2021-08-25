import { once, loadJs } from '/assets/js/util.js';

once(
    "js-downloader", 
    () => loadJs("/assets/js/lib/js-file-downloader.min.js")
);


