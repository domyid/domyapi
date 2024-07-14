const puppeteer = require('puppeteer');

async function getPdfUrl(fileName) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto('https://repo.ulbi.ac.id/buktiajar/#2023-2', { waitUntil: 'networkidle2' });

    const pdfUrl = await page.evaluate((fileName) => {
        const links = document.querySelectorAll('a');
        for (let link of links) {
            if (link.href.includes(fileName)) {
                return link.href;
            }
        }
        return null;
    }, fileName);

    await browser.close();
    return pdfUrl;
}

module.exports = { getPdfUrl };
