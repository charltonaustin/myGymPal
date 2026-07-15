/**
 * makeAutocomplete(input, names[, onSelect])
 *
 * Attaches a custom dropdown autocomplete to `input`.
 *   input    – HTMLInputElement (should have autocomplete="off")
 *   names    – string[] of suggestion values
 *   onSelect – optional callback(name); defaults to setting input.value
 *
 * Keyboard: ArrowDown/ArrowUp to move highlight, Enter to select, Escape to close.
 */
function makeAutocomplete(input, names, onSelect) {
    var wrapper = input.parentElement;
    wrapper.style.position = 'relative';
    var dropdown = null;
    var activeIndex = -1;

    function getItems() {
        return dropdown ? Array.from(dropdown.querySelectorAll('button')) : [];
    }

    function setActive(index) {
        var items = getItems();
        items.forEach(function (btn, i) {
            btn.classList.toggle('active', i === index);
        });
        activeIndex = index;
        if (index >= 0 && index < items.length) {
            items[index].scrollIntoView({ block: 'nearest' });
        }
    }

    function showDropdown(matches) {
        if (!dropdown) {
            dropdown = document.createElement('div');
            dropdown.className = 'list-group shadow';
            dropdown.style.cssText = 'position:absolute;top:100%;left:0;right:0;z-index:1050;max-height:220px;overflow-y:auto;border-radius:0 0 .375rem .375rem;';
            wrapper.appendChild(dropdown);
        }
        dropdown.innerHTML = '';
        activeIndex = -1;
        matches.forEach(function (name) {
            var btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'list-group-item list-group-item-action py-2 px-3';
            btn.style.fontSize = '0.95rem';
            btn.textContent = name;
            btn.addEventListener('mousedown', function (e) {
                e.preventDefault();
                selectName(name);
            });
            dropdown.appendChild(btn);
        });
    }

    function hideDropdown() {
        if (dropdown) { dropdown.remove(); dropdown = null; }
        activeIndex = -1;
    }

    function selectName(name) {
        input.value = name;
        if (typeof onSelect === 'function') onSelect(name);
        hideDropdown();
    }

    input.addEventListener('input', function () {
        var val = input.value.toLowerCase().trim();
        if (!val) { hideDropdown(); return; }
        var matches = names.filter(function (n) { return n.toLowerCase().includes(val); }).slice(0, 10);
        if (matches.length) showDropdown(matches); else hideDropdown();
    });

    input.addEventListener('keydown', function (e) {
        var items = getItems();
        if (e.key === 'ArrowDown') {
            e.preventDefault();
            setActive(Math.min(activeIndex + 1, items.length - 1));
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            setActive(Math.max(activeIndex - 1, -1));
        } else if (e.key === 'Enter') {
            if (dropdown && activeIndex >= 0 && activeIndex < items.length) {
                e.preventDefault();
                selectName(items[activeIndex].textContent);
            }
        } else if (e.key === 'Escape') {
            hideDropdown();
        }
    });

    input.addEventListener('blur', function () { setTimeout(hideDropdown, 150); });
}
