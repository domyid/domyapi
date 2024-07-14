const puppeteer = require('puppeteer');

async function getPdfUrl(fileName) {
    const browser = await puppeteer.launch({ headless: true });
    const page = await browser.newPage();
    await page.goto('https://repo.ulbi.ac.id/buktiajar/#2023-2');

    // Tunggu hingga halaman termuat
    await page.waitForTimeout(2000);

    const pdfURL = await page.evaluate((fileName) => {
        const links = Array.from(document.querySelectorAll('a'));
        const link = links.find(a => a.href.includes(fileName));
        return link ? link.href : null;
    }, fileName);

    await browser.close();

    if (!pdfURL) {
        throw new Error('Failed to find PDF URL on repository page');
    }

    return pdfURL;
}

// Contoh penggunaan fungsi
(async () => {
    try {
        const pdfURL = await getPdfUrl("BAP-Kecerdasan_Buatan_Artificial_Intelligence-51.pdf");
        console.log("PDF URL:", pdfURL);
    } catch (err) {
        console.error("Error:", err);
    }
})();
