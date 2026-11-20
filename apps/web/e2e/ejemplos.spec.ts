import { expect, test } from '@playwright/test'

test('los ejemplos de plantilla permiten crear y validar entradas con MSW', async ({
  page,
}) => {
  await page.goto('/?mocks=1')
  await expect(
    page.getByRole('heading', { name: 'Primer ejemplo', exact: true }),
  ).toBeVisible()
  await page.getByRole('button', { name: 'Crear ejemplo', exact: true }).click()
  await expect(page.getByText('El nombre es obligatorio.', { exact: true })).toBeVisible()
  await page.getByLabel('Nombre', { exact: true }).fill('Ejemplo de navegador')
  await page.getByLabel('Descripción', { exact: true }).fill('Prueba de la plantilla')
  await page.getByRole('button', { name: 'Crear ejemplo', exact: true }).click()
  await expect(
    page.getByRole('heading', { name: 'Ejemplo de navegador', exact: true }),
  ).toBeVisible()
  await expect(page.getByLabel('Nombre', { exact: true })).toHaveValue('')
})
