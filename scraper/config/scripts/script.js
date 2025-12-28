// that's the main script that should be injected and return data from the page.
// you can modify selectors here if you use different provider (or provided site changes classes).

const url = window.location.href;

if (url.includes('collection') || url.includes('cases')){
    Array.from(document.querySelectorAll('.item-box')).map(container => ({
        weapon: container.querySelector('.collection-skinbox span')?.textContent.trim() || '',
        name: container.querySelector('.block:not(.txt-small):not(.txt-stattrak)')?.textContent.trim() || '',
        rarity: container.querySelector('.classified, .covert, .restricted, .mil-spec, .industrial, .consumer')?.textContent.trim() || '',
        price: container.querySelector('.block.txt-small:not(.txt-stattrak):not(.txt-dark-grey)')?.textContent.trim() || '',
        stattrakPrice: container.querySelector('.block.txt-small.txt-stattrak, .block.txt-small.txt-souvenir')?.textContent.trim() || '',
        collection: container.querySelector('.block:not(.txt-small):not(.txt-stattrak):not(.txt-white)')?.textContent.trim() || '',
        url: container.querySelector('img')?.src?.replace('/webp/', '/').replace('.webp', '.png') || ''
     }));
} else if (url.includes('agents')){
    //FIX: side is not compatible. should work with selectors
    Array.from(document.querySelectorAll('.item-box')).map(container => ({
        name: container.querySelector('.block.txt-med.txt-white')?.textContent.trim() || '',
        affiliation: container.querySelector('.block.txt-small.txt-dark-grey')?.textContent.trim() || '',
        side: document.querySelectorAll('.flex-grow.txt-right')[2]?.textContent.trim() || '',
        collection: container.querySelector('.block:not(.txt-med):not(.txt-small):not(.txt-dark-grey)')?.textContent.trim() || '',
        rarity: container.querySelector('.classified, .covert, .restricted, .mil-spec, .industrial, .consumer')?.textContent.trim() || '',
        price: container.querySelector('.block.txt-small:not(.txt-dark-grey)')?.textContent.trim() || '',
        url: container.querySelector('img')?.src?.replace('/webp/', '/').replace('.webp', '.png') || ''
    }));
}
