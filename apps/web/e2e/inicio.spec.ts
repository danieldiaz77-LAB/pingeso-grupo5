import { expect, test } from '@playwright/test'
test('la aplicación monta el contenido en un navegador real', async ({ page }) => {
  await page.goto('/')
  await expect(
    page.getByRole('heading', { name: 'Entorno técnico PINGESO' }),
  ).toBeVisible()
  await expect(
    page.getByText('Contrato funcional pendiente', { exact: false }),
  ).toBeVisible()
})
