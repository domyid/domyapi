const functions = require('@google-cloud/functions-framework');
const puppeteer = require('puppeteer');

functions.http('getPdfUrl', async (req, res) => {
    const fileName = req.query.fileName;
    const url = 'https://repo.ulbi.ac.id/buktiajar/#2023-2';

    try {
        const browser = await puppeteer.launch({
            args: ['--no-sandbox', '--disable-setuid-sandbox'],
        });
        const page = await browser.newPage();
        await page.goto(url, { waitUntil: 'networkidle2' });
        await page.waitForTimeout(2000); // Menunggu halaman termuat

        const pdfURL = await page.evaluate((fileName) => {
            const links = Array.from(document.querySelectorAll('a'));
            const link = links.find(a => a.href.includes(fileName));
            return link ? link.href : null;
        }, fileName);

        await browser.close();

        if (pdfURL) {
            res.status(200).send(pdfURL);
        } else {
            res.status(404).send('File not found');
        }
    } catch (error) {
        res.status(500).send(`Error: ${error.message}`);
    }
});
