# G9 - SBOM + Trivy Scan

## Amac

Projeyi SBOM olarak belgelemek ve Trivy ile container/dependency CVE taramasi yapmak.

## Komutlar

Syft ile SBOM:

```bash
syft dir:. -o cyclonedx-json > gorev-09/sbom.json
```

Trivy filesystem taramasi:

```bash
trivy fs --severity CRITICAL,HIGH .
```

Trivy image taramasi:

```bash
docker compose build
trivy image --severity CRITICAL,HIGH secscan-backend:latest
trivy image --severity CRITICAL,HIGH secscan-frontend:latest
```

Docker ile Syft/Trivy alternatifi:

```bash
docker run --rm -v "${PWD}:/src" anchore/syft:latest dir:/src -o cyclonedx-json > gorev-09/sbom.json
docker run --rm -v "${PWD}:/src" aquasec/trivy:latest fs --severity CRITICAL,HIGH /src
```

## Uyguladigim Adimlar

1. Projenin SBOM dosyasini urettim.
2. Trivy ile critical/high bulgulari taradim.
3. Her CVE icin paket, seviye, fix version ve not bilgisi tabloya yazdim.
4. En az bir uygulanabilir fix varsa dependency versiyonunu guncelledim.

## Bulgular

| CVE ID | Package | Severity | Current | Fixed Version | Aksiyon |
|---|---|---|---|---|---|
| [DOLDUR] | [DOLDUR] | [DOLDUR] | [DOLDUR] | [DOLDUR] | [DOLDUR] |

## Karsilastigim Hatalar / Debug Notlari

- [DOLDUR] Syft/Trivy localde kurulu muydu, Docker image ile mi calistirdin?
- [DOLDUR] Trivy database download ne kadar surdu?
- [DOLDUR] Hangi bulgu uygulaman icin gercek risk olusturuyor?

## Sonuc

SBOM, projedeki paketlerin envanterini verir. Trivy bu envanter ve image katmanlari uzerinden bilinen CVE'leri bulur; fakat her CVE icin exploitability ayrica degerlendirilmelidir.

## Ogrendigim 3 Sey

- SBOM olmadan dependency riskini takip etmek zordur.
- CVE skoru yuksek olsa bile uygulamada kullanilmayan kod yolu riski azaltabilir.
- Fix version bilgisi varsa en temiz cozum dependency guncellemesidir.
