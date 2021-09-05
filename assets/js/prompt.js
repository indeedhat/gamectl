import { Modal } from '/assets/js/modal.js';

/**
 * Emulate the window.confirm functionality with custom ui
 * 
 * @async
 *
 * @param string title
 * @param string body
 *
 * @return Promise<void>
 */
const confirm = async (title, body) => {
    let outcome = false;

    return new Promise((accept, reject) => {
        const modal = new Modal({
            title,
            content: _confirmTemplate(body),
            onClose: () => outcome ? accept() : reject(),
            onOpen: () => {
                const $form = modal.$body.querySelector("form");
                const $confirm = modal.$body.querySelector("button.prompt-confirm");
                const $cancel = modal.$body.querySelector("button.prompt-cancel");

                $form.onsubmit = e => {
                    e.preventDefault();
                };

                $confirm.onclick = e => {
                    e.preventDefault();
                    e.stopPropagation();
                    outcome = true;
                    modal.close();
                };

                $cancel.onclick = e => {
                    e.preventDefault();
                    e.stopPropagation();
                    modal.close();
                };
            }
        });

        modal.open();
    });
};

const _confirmTemplate = body => {
    return `
        <div>${body}</div>
        <form>
            <div class="group">
                <button class="prompt-confirm">Confirm</button>
                <button class="prompt-cancel">Cancel</button>
            </div>
        </form>
    `;
};

export {
    confirm
};
