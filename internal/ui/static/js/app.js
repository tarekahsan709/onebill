// Sentinel values for specific list navigation.

// Simple Polyfill for browsers that don't support Trusted Types
// See https://caniuse.com/?search=trusted%20types
if (!window.trustedTypes || !trustedTypes.createPolicy) {
    window.trustedTypes = {
        createPolicy: (name, policy) => ({
            createScriptURL: src => src,
            createHTML: html => html,
        })
    };
}

/**
 * Send a POST request to the specified URL with the given body.
 *
 * @param {string} url - The URL to send the request to.
 * @param {Object} [body] - The body of the request (optional).
 * @returns {Promise<Response>} The response from the fetch request.
 */
function sendPOSTRequest(url, body = null) {
    const options = {
        method: "POST",
        headers: {
            "X-Csrf-Token": document.body.dataset.csrfToken || ""
        }
    };

    if (body !== null) {
        options.headers["Content-Type"] = "application/json";
        options.body = JSON.stringify(body);
    }

    return fetch(url, options);
}

/**
 * Attach a click event listener to elements matching the selector.
 *
 * @param {string} selector
 * @param {function} callback
 * @param {boolean} noPreventDefault
 */
function onClick(selector, callback, noPreventDefault) {
    document.querySelectorAll(selector).forEach((element) => {
        element.onclick = (event) => {
            if (!noPreventDefault) {
                event.preventDefault();
            }
            callback(event);
        };
    });
}

/**
 * Filter visible elements based on the selector.
 *
 * @param {string} selector
 * @returns {Array<Element>}
 */
function getVisibleElements(selector) {
    const elements = document.querySelectorAll(selector);
    return [...elements].filter((element) => element.offsetParent !== null);
}

/**
 * Get all visible entries on the current page.
 *
 * @return {Array<Element>}
 */
function getVisibleEntries() {
    return getVisibleElements(".items .item");
}

/**
 * Check if the current view is a list view.
 *
 * @returns {boolean}
 */
function isListView() {
    return document.querySelector(".items") !== null;
}

/**
 * Find the entry element for the given element.
 *
 * @returns {Element|null}
 */
function findEntry(element) {
    if (element) {
        return element.closest(".item");
    }
    return document.querySelector(".current-item");
}

/**
 * Create an icon label element with the given text.
 *
 * @param {string} labelText - The text to display in the icon label.
 * @returns {Element} The created icon label element.
 */
function createIconLabelElement(labelText) {
    const labelElement = document.createElement("span");
    labelElement.classList.add("icon-label");
    labelElement.textContent = labelText;
    return labelElement;
}

/**
 * Set the icon and label element in the parent element.
 *
 * @param {Element} parentElement - The parent element to insert the icon and label into.
 * @param {string} iconName - The name of the icon to display.
 * @param {string} labelText - The text to display in the label.
 */
function setIconAndLabelElement(parentElement, iconName, labelText) {
    const iconElement = document.querySelector(`template#icon-${iconName}`);
    if (iconElement) {
        const iconClone = iconElement.content.cloneNode(true);
        parentElement.textContent = ""; // Clear existing content
        parentElement.appendChild(iconClone);
    }

    if (labelText) {
        const labelElement = createIconLabelElement(labelText);
        parentElement.appendChild(labelElement);
    }
}

/**
 * Set the button to a loading state and return a clone of the original button element.
 *
 * @param {Element} buttonElement - The button element to set to loading state.
 * @return {Element} The original button element cloned before modification.
 */
function setButtonToLoadingState(buttonElement) {
    const originalButtonElement = buttonElement.cloneNode(true);

    buttonElement.textContent = "";
    buttonElement.appendChild(createIconLabelElement(buttonElement.dataset.labelLoading));

    return originalButtonElement;
}

/**
 * Set the read status button state.
 *
 * @param {Element} buttonElement - The button element to update.
 * @param {string} newState - The new state to set ("read" or "unread").
 */
function setReadStatusButtonState(buttonElement, newState) {
    buttonElement.dataset.value = newState;
    const iconType = newState === "read" ? "unread" : "read";
    setIconAndLabelElement(buttonElement, iconType, buttonElement.dataset[newState === "read" ? "labelUnread" : "labelRead"]);
}

/**
 * Navigate to a specific page.
 *
 * @param {string} page - The page to navigate to.
 * @param {boolean} reloadOnFail - If true, reload the current page if the target page is not found.
 */
function goToPage(page, reloadOnFail = false) {
    const element = document.querySelector(`:is(a, button)[data-page=${page}]`);

    if (element) {
        document.location.href = element.href;
    } else if (reloadOnFail) {
        window.location.reload();
    }
}

/**
 * Toggle the ARIA attributes on the main menu based on the viewport width.
 */
function toggleAriaAttributesOnMainMenu() {
    const logoElement = document.querySelector(".logo");
    const homePageLinkElement = document.querySelector(".logo > a");

    if (!logoElement || !homePageLinkElement) return;

    const isMobile = document.documentElement.clientWidth < 650;

    if (isMobile) {
        const navMenuElement = document.getElementById("header-menu");
        const isExpanded = navMenuElement?.classList.contains("js-menu-show") ?? false;
        const toggleButtonLabel = logoElement.getAttribute("data-toggle-button-label");

        // Set mobile menu button attributes
        Object.assign(logoElement, {
            role: "button",
            tabIndex: 0,
            ariaLabel: toggleButtonLabel,
            ariaExpanded: isExpanded.toString()
        });
        homePageLinkElement.tabIndex = -1;
    } else {
        // Remove mobile menu button attributes
        ["role", "tabindex", "aria-expanded", "aria-label"].forEach(attr => {
            logoElement.removeAttribute(attr);
        });
        homePageLinkElement.removeAttribute("tabindex");
    }
}

/**
 * Toggle the main menu dropdown.
 *
 * @param {Event} event - The event object.
 */
function toggleMainMenuDropdown(event) {
    // Only handle Enter, Space, or click events
    if (event.type === "keydown" && !["Enter", " "].includes(event.key)) {
        return;
    }

    // Prevent default only if element has role attribute (mobile menu button)
    if (event.currentTarget.getAttribute("role")) {
        event.preventDefault();
    }

    const navigationMenu = document.querySelector(".header nav ul");
    const menuToggleButton = document.querySelector(".logo");

    if (!navigationMenu || !menuToggleButton) {
        return;
    }

    const isShowing = navigationMenu.classList.toggle("js-menu-show");
    menuToggleButton.setAttribute("aria-expanded", isShowing.toString());
}

/**
 * Initialize the main menu handlers.
 */
function initializeMainMenuHandlers() {
    toggleAriaAttributesOnMainMenu();
    window.addEventListener("resize", toggleAriaAttributesOnMainMenu, { passive: true });

    const logoElement = document.querySelector(".logo");
    if (logoElement) {
        logoElement.addEventListener("click", toggleMainMenuDropdown);
        logoElement.addEventListener("keydown", toggleMainMenuDropdown);
    }

    onClick(".header nav li", (event) => {
        const linkElement = event.target.closest("a") || event.target.querySelector("a");
        if (linkElement && !event.ctrlKey && !event.shiftKey && !event.metaKey) {
            event.preventDefault();
            window.location.href = linkElement.getAttribute("href");
        }
    }, true);
}

/**
 * This function changes the button label to the loading state and disables the button.
 *
 * @returns {void}
 */
function initializeFormHandlers() {
    document.querySelectorAll("form").forEach((element) => {
        element.onsubmit = () => {
            const buttons = element.querySelectorAll("button[type=submit]");
            buttons.forEach((button) => {
                if (button.dataset.labelLoading) {
                    button.textContent = button.dataset.labelLoading;
                }
                button.disabled = true;
            });
        };
    });
}

/**
 * Mark all visible entries on the current page as read.
 */
function markPageAsReadAction() {
    const items = getVisibleEntries();
    if (items.length === 0) return;

    const entryIDs = items.map((element) => parseInt(element.dataset.id, 10));

    // Batch DOM writes after all reads
    items.forEach((element) => element.classList.add("item-status-read"));

    updateEntriesStatus(entryIDs, "read", () => {
        const element = document.querySelector(":is(a, button)[data-action=markPageAsRead]");
        const showOnlyUnread = element?.dataset.showOnlyUnread || false;

        if (showOnlyUnread) {
            window.location.reload();
        } else {
            goToPage("next", true);
        }
    });
}

/**
 * Toggle the read status of the entry that triggered the action.
 *
 * @param {Element} element Element that triggered the action.
 * @returns {void}
 */
function handleEntryStatus(element) {
    const currentEntry = findEntry(element);
    if (currentEntry) {
        toggleEntryStatus(currentEntry);
    }
}

/**
 * Toggle the entry status between "read" and "unread".
 *
 * @param {Element} element The entry element to toggle the status for.
 */
function toggleEntryStatus(element) {
    const entryID = parseInt(element.dataset.id, 10);
    const buttonElement = element.querySelector(":is(a, button)[data-toggle-status]");
    if (!buttonElement) return;

    const currentStatus = buttonElement.dataset.value;
    const newStatus = currentStatus === "read" ? "unread" : "read";

    setButtonToLoadingState(buttonElement);

    updateEntriesStatus([entryID], newStatus, () => {
        setReadStatusButtonState(buttonElement, newStatus);

        element.classList.replace(`item-status-${currentStatus}`, `item-status-${newStatus}`);

        if (isListView() && getVisibleEntries().length === 0) {
            window.location.reload();
        }
    });
}

/**
 * Update the status of multiple entries.
 *
 * @param {Array<number>} entryIDs - The IDs of the entries to update.
 * @param {string} status - The new status to set for the entries (e.g., "read", "unread").
 */
function updateEntriesStatus(entryIDs, status, callback) {
    const url = document.body.dataset.entriesStatusUrl;
    sendPOSTRequest(url, { entry_ids: entryIDs, status: status }).then((resp) => {
        resp.json().then(count => {
            if (callback) {
                callback(resp);
            }
            updateUnreadCounterValue(status === "read" ? -count : count);
        });
    });
}

/**
 * Update the unread counter value.
 *
 * @param {number} delta - The amount to change the counter by.
 */
function updateUnreadCounterValue(delta) {
    document.querySelectorAll("span.unread-counter").forEach((element) => {
        const oldValue = parseInt(element.textContent, 10);
        element.textContent = oldValue + delta;
    });

    if (window.location.href.endsWith('/unread')) {
        const oldValue = parseInt(document.title.split('(')[1], 10);
        const newValue = oldValue + delta;
        document.title = document.title.replace(/\(\d+\)/, `(${newValue})`);
    }
}

/**
 * Handle confirmation messages for actions that require user confirmation.
 *
 * This function modifies the link element to show a confirmation question with "Yes" and "No" buttons.
 * If the user clicks "Yes", it calls the provided callback with the URL and redirect URL.
 * If the user clicks "No", it either redirects to a no-action URL or restores the link element.
 *
 * @param {Element} linkElement - The link or button element that triggered the confirmation.
 * @param {function} callback - The callback function to execute if the user confirms the action.
 * @returns {void}
 */
function handleConfirmationMessage(linkElement, callback) {
    if (linkElement.tagName !== 'A' && linkElement.tagName !== "BUTTON") {
        linkElement = linkElement.parentNode;
    }

    linkElement.style.display = "none";

    const containerElement = linkElement.parentNode;
    const questionElement = document.createElement("span");

    function createLoadingElement() {
        const loadingElement = document.createElement("span");
        loadingElement.className = "loading";
        loadingElement.appendChild(document.createTextNode(linkElement.dataset.labelLoading));

        questionElement.remove();
        containerElement.appendChild(loadingElement);
    }

    const yesElement = document.createElement("button");
    yesElement.appendChild(document.createTextNode(linkElement.dataset.labelYes));
    yesElement.onclick = (event) => {
        event.preventDefault();

        createLoadingElement();

        callback(linkElement.dataset.url, linkElement.dataset.redirectUrl);
    };

    const noElement = document.createElement("button");
    noElement.appendChild(document.createTextNode(linkElement.dataset.labelNo));
    noElement.onclick = (event) => {
        event.preventDefault();

        const noActionUrl = linkElement.dataset.noActionUrl;
        if (noActionUrl) {
            createLoadingElement();

            callback(noActionUrl, linkElement.dataset.redirectUrl);
        } else {
            linkElement.style.display = "inline";
            questionElement.remove();
        }
    };

    questionElement.className = "confirm";
    questionElement.appendChild(document.createTextNode(`${linkElement.dataset.labelQuestion} `));
    questionElement.appendChild(yesElement);
    questionElement.appendChild(document.createTextNode(", "));
    questionElement.appendChild(noElement);

    containerElement.appendChild(questionElement);
}

/**
 * Initialize the service worker and PWA installation prompt.
 */
function initializeServiceWorker() {
    // Register service worker if supported
    if ("serviceWorker" in navigator) {
        const serviceWorkerURL = document.body.dataset.serviceWorkerUrl;
        if (serviceWorkerURL) {
            const ttpolicy = trustedTypes.createPolicy('url', {createScriptURL: src => src});
            navigator.serviceWorker.register(ttpolicy.createScriptURL(serviceWorkerURL), {
                type: "module"
            }).catch((error) => {
                console.error("Service Worker registration failed:", error);
            });
        }
    }

    // PWA installation prompt handling
    window.addEventListener("beforeinstallprompt", (event) => {
        let deferredPrompt = event;
        const promptHomeScreen = document.getElementById("prompt-home-screen");
        const btnAddToHomeScreen = document.getElementById("btn-add-to-home-screen");

        if (!promptHomeScreen || !btnAddToHomeScreen) return;

        promptHomeScreen.style.display = "block";

        btnAddToHomeScreen.addEventListener("click", (event) => {
            event.preventDefault();
            deferredPrompt.prompt();
            deferredPrompt.userChoice.then(() => {
                deferredPrompt = null;
                promptHomeScreen.style.display = "none";
            });
        });
    });
}

/**
 * Initialize WebAuthn handlers if supported.
 */
function initializeWebAuthn() {
    if (typeof WebAuthnHandler !== 'function') return;
    if (!WebAuthnHandler.isWebAuthnSupported()) return;

    const webauthnHandler = new WebAuthnHandler();

    // Setup delete credentials handler
    onClick("#webauthn-delete", () => { webauthnHandler.removeAllCredentials(); });

    // Setup registration
    const registerButton = document.getElementById("webauthn-register");
    if (registerButton) {
        registerButton.disabled = false;
        onClick("#webauthn-register", () => {
            webauthnHandler.register().catch((err) => WebAuthnHandler.showErrorMessage(err));
        });
    }

    // Setup login
    const loginButton = document.getElementById("webauthn-login");
    const usernameField = document.getElementById("form-username");

    if (loginButton && usernameField) {
        const abortController = new AbortController();
        loginButton.disabled = false;

        onClick("#webauthn-login", () => {
            abortController.abort();
            webauthnHandler.login().catch(err => WebAuthnHandler.showErrorMessage(err));
        });

        webauthnHandler.conditionalLogin(abortController).catch(err => WebAuthnHandler.showErrorMessage(err));
    }
}

/**
 * Initialize click handlers for various UI elements.
 */
function initializeClickHandlers() {
    // Entry actions
    onClick(":is(a, button)[data-toggle-status]", (event) => handleEntryStatus(event.target));

    // Page actions with confirmation
    onClick(":is(a, button)[data-action=markPageAsRead]", (event) => handleConfirmationMessage(event.target, markPageAsReadAction));

    // Generic confirmation handler
    onClick(":is(a, button)[data-confirm]", (event) => {
        handleConfirmationMessage(event.target, (url, redirectURL) => {
            sendPOSTRequest(url).then((response) => {
                if (redirectURL) {
                    window.location.href = redirectURL;
                } else if (response?.redirected && response.url) {
                    window.location.href = response.url;
                } else {
                    window.location.reload();
                }
            });
        });
    });

}

// Initialize application handlers
initializeMainMenuHandlers();
initializeFormHandlers();
initializeWebAuthn();
initializeClickHandlers();
initializeServiceWorker();
