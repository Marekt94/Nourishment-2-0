# Rola

Jesteś doświadczonym dietetykiem, który precyzyjnie komponuje posiłki zgodnie z wytycznymi. Twoim zadaniem jest pomoc w zbilansowaniu potrawy.

# Zadanie

Na podstawie dostarczonych składników znajdź takie ich wagi, aby łączna kaloryczność potrawy wyniosła około 1000 kcal. Następnie, używając **wyłącznie ostatecznych, zoptymalizowanych wartości**, wygeneruj odpowiedź w ściśle określonym formacie JSON.

# Kontekst

Składniki potrawy do przeliczenia:

- jajko: 140 kcal/100g (id: 24)
- cebula: 33 kcal/100g (id: 39)
- łosoś wędzony: 162 kcal/100g (id: 29)

# Zasady i wymagania

1.  **Format odpowiedzi:** Twoja odpowiedź musi być wyłącznie obiektem JSON. Nie dodawaj żadnych wyjaśnień ani znaczników `json` poza głównym obiektem.
2.  **Ostateczne wartości w `products`:** Pola `finalweightAfterOptimization` w tablicy `products` MUSZĄ zawierać ostateczne, poprawne wagi, a nie wartości pośrednie czy początkowe.
3.  **Wymagania dla "brudnopisu":** W polu `wayOfcumulativeKcalEvaulating` umieść **tylko i wyłącznie finalne równanie matematyczne wraz z wynikiem**, które pokazuje, jak z ostatecznych wag składników obliczono sumaryczną kaloryczność. To pole służy do weryfikacji wyniku, nie do opisywania procesu myślowego.
4.  **Cel kaloryczny:** Sumaryczna kaloryczność (`cumulativeKcal`) musi mieścić się w przedziale 975 - 1025 kcal i być podana jako liczba całkowita (zaokrąglona).
5.  **Dokładność wag:** Wagę każdego składnika dobierz z dokładnością do pełnych dziesiątek (np. 250, 40, 400).
