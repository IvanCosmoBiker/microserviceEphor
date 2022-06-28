package payments

import (
    sberBank "ephorservices/src/server/services/payments/sber"
    vendpay "ephorservices/src/server/services/payments/vendpay"
    interfaceBank "ephorservices/src/server/utils/interface/payment"
)
// instance of type banks
var bankSber sberBank.NewSberStruct
var bankVendPay vendpay.NewVendStruct

func GetBank(bank int) (interfaceBank.Bank) {
  switch bank {
        case interfaceBank.TypeSber:
            return bankSber.NewBank()
            fallthrough
        case interfaceBank.TypeVendPay:
            return bankVendPay.NewBank()
  }
  return nil
}



