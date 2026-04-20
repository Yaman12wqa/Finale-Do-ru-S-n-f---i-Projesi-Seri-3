# G8 - SAST Pipeline (Semgrep)

## Amac

GitHub Actions uzerinden Semgrep SAST kontrolu calistirmak ve bulgulari raporlamak.

## Repo Durumu

Bu projede Semgrep zaten `.github/workflows/security.yml` icinde bulunuyor. Workflow su kurallari calistirir:

- `p/owasp-top-ten`
- `p/golang`
- `p/typescript`

## Uyguladigim Adimlar

1. `.github/workflows/security.yml` dosyasini kontrol ettim.
2. Repoyu GitHub'a push edecegim.
3. Actions sekmesinde Security workflow'unun calistigini dogrulayacagim.
4. Bilerek guvensiz kod ekleme deneyi icin ayri bir lab branch kullanacagim.
5. Semgrep bulgusunu ekran goruntusu ile rapora ekleyecegim.

## Guvensiz Kod Deneyi

Ana projeye kalici zafiyet eklemek dogru degildir. Deney icin ayri branch ac:

```bash
git checkout -b lab/semgrep-demo
```

Ornek olarak gecici bir dosyada su tarz bir kod denenebilir:

```js
app.get("/debug", (req, res) => {
  eval(req.query.code);
  res.send("ok");
});
```

Semgrep bulgusu alindiktan sonra bu branch merge edilmemelidir.

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] GitHub Actions ilk calistiginda hangi sonuc geldi?
- [DOLDUR] Semgrep hangi rule id ile bulgu verdi?
- [DOLDUR] Bulguyu nasil duzelttin veya neden lab branchte biraktin?

## Sonuc

SAST pipeline, kod GitHub'a push edildiginde otomatik guvenlik kontrolu saglar. Bilerek zafiyet ekleme testi ana branchte kalici tutulmamalidir.

## Ogrendigim 3 Sey

- Semgrep pattern tabanli calisir ve tehlikeli API kullanimlarini erken yakalar.
- CI/CD guvenlik kontrolu manuel review yerine gecmez, ama iyi bir erken uyari sistemidir.
- Guvensiz demo kodlari production branchte kalmamali, lab branchte tutulmalidir.
