import React, { FC, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import type * as Type from '@/common/interface';
import { useToast } from '@/hooks';
import { useSmtpSetting, updateSmtpSetting } from '@/services';
import pattern from '@/common/pattern';
import { SchemaForm, JSONSchema, UISchema } from '@/components';
import { initFormData } from '../../../components/SchemaForm/index';
import { handleFormError } from '@/utils';

const Smtp: FC = () => {
  const { t } = useTranslation('translation', {
    keyPrefix: 'admin.smtp',
  });
  const Toast = useToast();
  const { data: setting } = useSmtpSetting();
  const schema: JSONSchema = {
    title: t('page_title'),
    properties: {
      from_email: {
        type: 'string',
        title: t('from_email.label'),
        description: t('from_email.text'),
      },
      from_name: {
        type: 'string',
        title: t('from_name.label'),
        description: t('from_name.text'),
      },
      smtp_host: {
        type: 'string',
        title: t('smtp_host.label'),
        description: t('smtp_host.text'),
      },
      encryption: {
        type: 'boolean',
        title: t('encryption.label'),
        description: t('encryption.text'),
        enum: ['SSL', ''],
        enumNames: ['SSL', 'None'],
      },
      smtp_port: {
        type: 'string',
        title: t('smtp_port.label'),
        description: t('smtp_port.text'),
      },
      smtp_authentication: {
        type: 'boolean',
        title: t('smtp_authentication.label'),
        enum: [true, false],
        enumNames: [t('smtp_authentication.yes'), t('smtp_authentication.no')],
      },
      smtp_username: {
        type: 'string',
        title: t('smtp_username.label'),
      },
      smtp_password: {
        type: 'string',
        title: t('smtp_password.label'),
      },
      test_email_recipient: {
        type: 'string',
        title: t('test_email_recipient.label'),
        description: t('test_email_recipient.text'),
      },
    },
  };
  const uiSchema: UISchema = {
    from_email: {
      'ui:options': {
        type: 'email',
      },
    },
    encryption: {
      'ui:widget': 'select',
    },
    smtp_username: {
      'ui:options': {
        validator: (value: string, formData) => {
          if (formData.smtp_authentication.value) {
            if (!value) {
              return t('smtp_username.msg');
            }
          }
          return true;
        },
      },
    },
    smtp_password: {
      'ui:options': {
        type: 'password',
        validator: (value: string, formData) => {
          if (formData.smtp_authentication.value) {
            if (!value) {
              return t('smtp_password.msg');
            }
          }
          return true;
        },
      },
    },
    smtp_authentication: {
      'ui:widget': 'switch',
    },
    smtp_port: {
      'ui:options': {
        type: 'number',
        validator: (value) => {
          if (!/^[1-9][0-9]*$/.test(value) || Number(value) > 65535) {
            return t('smtp_port.msg');
          }
          return true;
        },
      },
    },
    test_email_recipient: {
      'ui:options': {
        type: 'email',
        validator: (value) => {
          if (value && !pattern.email.test(value)) {
            return t('test_email_recipient.msg');
          }
          return true;
        },
      },
    },
  };
  const [formData, setFormData] = useState<Type.FormDataType>(
    initFormData(schema),
  );

  const onSubmit = (evt) => {
    evt.preventDefault();
    evt.stopPropagation();

    const reqParams: Type.AdminSettingsSmtp = {
      from_email: formData.from_email.value,
      from_name: formData.from_name.value,
      smtp_host: formData.smtp_host.value,
      encryption: formData.encryption.value,
      smtp_port: Number(formData.smtp_port.value),
      smtp_authentication: formData.smtp_authentication.value,
      ...(formData.smtp_authentication.value
        ? { smtp_username: formData.smtp_username.value }
        : {}),
      ...(formData.smtp_authentication.value
        ? { smtp_password: formData.smtp_password.value }
        : {}),
      test_email_recipient: formData.test_email_recipient.value,
    };

    updateSmtpSetting(reqParams)
      .then(() => {
        Toast.onShow({
          msg: t('update', { keyPrefix: 'toast' }),
          variant: 'success',
        });
      })
      .catch((err) => {
        if (err.isError) {
          const data = handleFormError(err, formData);
          setFormData({ ...data });
        }
      });
  };

  useEffect(() => {
    if (!setting) {
      return;
    }
    const formMeta = {};
    Object.keys(setting).forEach((k) => {
      formMeta[k] = { ...formData[k], value: setting[k] };
    });
    setFormData({ ...formData, ...formMeta });
  }, [setting]);

  useEffect(() => {
    if (formData.smtp_authentication.value === '') {
      return;
    }
    if (formData.smtp_authentication.value) {
      setFormData({
        ...formData,
        smtp_username: { ...formData.smtp_username, hidden: false },
        smtp_password: { ...formData.smtp_password, hidden: false },
      });
    } else {
      setFormData({
        ...formData,
        smtp_username: { ...formData.smtp_username, hidden: true },
        smtp_password: { ...formData.smtp_password, hidden: true },
      });
    }
  }, [formData.smtp_authentication.value]);

  const handleOnChange = (data) => {
    setFormData(data);
  };
  return (
    <>
      <h3 className="mb-4">{t('page_title')}</h3>
      <SchemaForm
        schema={schema}
        uiSchema={uiSchema}
        formData={formData}
        onChange={handleOnChange}
        onSubmit={onSubmit}
      />
    </>
  );
};

export default Smtp;
