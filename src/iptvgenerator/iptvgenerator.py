import random
import string

import requests

SIGNUP_URL = "https://server.azsat.org/novo/usuario/registrar_iptv.php"
DASH_URL = "https://server.azsat.org/"

def main():
    s = requests.Session()
    response = s.get(SIGNUP_URL)

    start = response.text.find('<input name="csrf" type=\'hidden\' value="') + 40
    end = response.text.find('"/>\n            <input name="fg" id="fg" type="hidden"/>')

    csrf = response.text[start:end]

    user = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(10))

    headers = {
        'Content-Type': "application/x-www-form-urlencoded"
    }
    payload = {
        "step": "2",
        "info": "dXNlcl9hZ2VudD0hPU1vemlsbGEvNS4wIChYMTE7IExpbnV4IHg4Nl82NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzcyLjAuMzYyNi45NiBTYWZhcmkvNTM3LjM2QEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxsYW5ndWFnZT0hPXB0LUJSQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxjb2xvcl9kZXB0aD0hPTI0QEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxwaXhlbF9yYXRpbz0hPTFAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB8fHx8fHx8fHx8fHx8fGhhcmR3YXJlX2NvbmN1cnJlbmN5PSE9NEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8cmVzb2x1dGlvbj0hPTE5MjAsMTA4MEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8YXZhaWxhYmxlX3Jlc29sdXRpb249IT0xOTIwLDEwNDBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB8fHx8fHx8fHx8fHx8fHRpbWV6b25lX29mZnNldD0hPTEyMEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8c2Vzc2lvbl9zdG9yYWdlPSE9MUBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8bG9jYWxfc3RvcmFnZT0hPTFAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB8fHx8fHx8fHx8fHx8fGluZGV4ZWRfZGI9IT0xQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxvcGVuX2RhdGFiYXNlPSE9MUBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8Y3B1X2NsYXNzPSE9dW5rbm93bkBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8bmF2aWdhdG9yX3BsYXRmb3JtPSE9TGludXggeDg2XzY0QEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxkb19ub3RfdHJhY2s9IT0xQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxyZWd1bGFyX3BsdWdpbnM9IT1DaHJvbWUgUERGIFBsdWdpbjo6UG9ydGFibGUgRG9jdW1lbnQgRm9ybWF0OjphcHBsaWNhdGlvbi94LWdvb2dsZS1jaHJvbWUtcGRmfnBkZixDaHJvbWUgUERGIFZpZXdlcjo6OjphcHBsaWNhdGlvbi9wZGZ",
        "fg": "ab2fa79efa2c077b8ab7232cc2646770",
        "csrf": csrf,
        "login": user,
        "email": user+"@gmail.com",
        "senha": user,
        "senha2": user,
        "operadora[]": "IPTV",
        "operadora[]": "IPTV",
        "operadora[]": "IPTV"
    }
    response = s.post(SIGNUP_URL, data=payload, headers=headers)

    response = s.get(DASH_URL)

    response = s.get(DASH_URL)

    start = response.text.find('<input name="url" readonly onclick="this.select();" value="') + 59
    end = response.text.find('" size="32"></h6>\n')

    url = response.text[start:end]

    file = open("/data/playlist-url", "w+")
    file.write(url)
    file.close()

    print(url)


if __name__ == "__main__":
    main()
