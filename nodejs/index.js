const fs = require('fs');
const puppeteer = require('puppeteer');

async function getPdfUrl(fileName) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto('https://repo.ulbi.ac.id/buktiajar/#2023-2');
    await page.waitForTimeout(2000); // Tambahkan delay untuk memastikan halaman termuat
    const content = await page.content();

    const startIndex = content.indexOf(fileName);
    if (startIndex === -1) {
        await browser.close();
        console.log('file not found');
        return;
    }

    const hrefStart = content.lastIndexOf('href="', startIndex) + 6;
    const hrefEnd = content.indexOf('"', hrefStart);
    const pdfURL = content.substring(hrefStart, hrefEnd);

    await browser.close();
    console.log(pdfURL);
}

const fileName = process.argv[2];
getPdfUrl(fileName);
