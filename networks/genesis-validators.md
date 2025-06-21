
---

# 🌐 Genesis Validator Setup Guide — **Symphony Mainnet**

> 🧾 *This guide will help you set up a Symphony genesis validator, generate your gentx with self-delegation, and submit a Pull Request (PR) to the official repository.*

---

## ✅ Minimum Server Requirements

| Component | Minimum              |
| --------- | -------------------- |
| OS        | Ubuntu 20.04 / 22.04 |
| CPU       | 4 Physical Cores     |
| RAM       | 16 GB                |
| Storage   | 500 GB SSD           |
| Network   | 100 Mbps             |

---

## 🔧 Step 1: Install Go (v1.16+)

```bash
sudo rm -rf /usr/local/go
curl -s https://raw.githubusercontent.com/canha/golang-tools-install-script/master/goinstall.sh | bash
source ~/.profile
go version
```

---

## 🏗️ Step 2: Clone and Build `symphonyd`

```bash
cd $HOME
git clone https://github.com/symphony-protocol/symphony
cd symphony
git checkout v1.0.0
make install
```

Verify:

```bash
symphonyd version --long
```

---

## ⚙️ Step 3: Initialize Node

```bash
symphonyd init <YOUR_MONIKER> --chain-id symphony-1
```

> Replace `<YOUR_MONIKER>` with your validator name (e.g. `mynode`)

---

## 📥 Step 4: Download Genesis File

```bash
curl -s https://raw.githubusercontent.com/Orchestra-Labs/symphony/refs/heads/main/networks/symphony-1/genesis.json \
-o ~/.symphonyd/config/genesis.json
```

---

## 🔐 Step 5: Validate Genesis File

```bash
symphonyd validate-genesis
```

---

## 🔑 Step 6: Import Wallet from Testnet
```bash
symphonyd keys add wallet --recover
```

## 📝 Step 7: Create `gentx` 

```bash
symphonyd gentx wallet 1000000note \
  --chain-id symphony-1 \
  --moniker "<YOUR_MONIKER>" \
  --identity "<YOUR_KEYBASE_ID>" \
  --website "https://<YOUR_WEBSITE>" \
  --security-contact "<YOUR_EMAIL>" \
  --details "<BRIEF_DESCRIPTION_OF_YOUR_VALIDATOR>" \
  --commission-rate 0.05 \
  --commission-max-rate 0.20 \
  --commission-max-change-rate 0.05 \
  --min-self-delegation "1"
```

> Replace values with your own info.

---

### Submit Your GenTx


## 📁 Step 1: Fork the Official Symphony Repo
> Fork https://github.com/Orchestra-Labs/symphony to your GitHub, then clone it:


```bash
git clone https://github.com/<YOUR_GITHUB_USERNAME>/symphony
cd symphony
mkdir -p networks/symphony-1/gentxs
cp ~/.symphonyd/config/gentx/*.json networks/symphony-1/gentxs/<your_moniker>.json
```

Check the file:

```bash
cat networks/symphony-1/gentxs/<your_moniker>.json | jq
```

---

## 🚀 Step 2: Submit a Pull Request (PR)

### 2.1 Commit and Push the Gentx File

```bash
git add networks/symphony-1/gentxs/<your_moniker>.json
git commit -m "Add gentx for <YOUR_MONIKER>"
git push origin main
```

> 💡 **Authentication Notice:**
> If you're prompted for a GitHub username and password when pushing:
>
> * Enter your **GitHub username**
> * For the password, **do not use your GitHub login password.** Instead, use a [Personal Access Token (Classic)](https://github.com/settings/tokens) with appropriate scopes (e.g. `repo`)
>
> This is required because GitHub no longer supports password-based authentication over HTTPS.

---

### 2.2 Open Pull Request

1. Go to your fork: `https://github.com/<your_github_username>/symphony`
2. Click **“Compare & pull request”**
3. Ensure base repo is `Orchestra-Labs/symphony` and target branch is `main`
4. Title: `Create gentx-<YOUR_MONIKER>.json`
5. Click **“Create Pull Request”**

---

## 🏁 Done!

## 🧩 Example gentx values:

| Field            | Example                                 |
| ---------------- | --------------------------------------- |
| moniker          | `"Sychonix"`                            |
| identity         | `"6366E6C36DFCFCA7"` (Keybase ID)       |
| website          | `"https://sychonix.com"`                |
| security-contact | `"yourname@domain.com"`                 |
| details          | `"Reliable validator and Web3 builder"` |
| amount           | `1000000note`                           |

---
