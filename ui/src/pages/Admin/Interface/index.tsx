import { FC, FormEvent, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { useToast } from '@/hooks';
import {
  LangsType,
  FormDataType,
  AdminSettingsInterface,
} from '@/common/interface';
import { interfaceStore } from '@/stores';
import { JSONSchema, SchemaForm, UISchema } from '@/components';
import { DEFAULT_TIMEZONE } from '@/common/constants';
import {
  updateInterfaceSetting,
  useInterfaceSetting,
  useThemeOptions,
} from '@/services';
import {
  setupAppLanguage,
  loadLanguageOptions,
  setupAppTimeZone,
} from '@/utils/localize';
import { handleFormError } from '@/utils';

const Interface: FC = () => {
  const { t } = useTranslation('translation', {
    keyPrefix: 'admin.interface',
  });
  const storeInterface = interfaceStore.getState().interface;
  const { data: themes } = useThemeOptions();
  const Toast = useToast();
  const [langs, setLangs] = useState<LangsType[]>();
  const { data: setting } = useInterfaceSetting();

  const schema: JSONSchema = {
    title: t('page_title'),
    properties: {
      theme: {
        type: 'string',
        title: t('theme.label'),
        description: t('theme.text'),
        enum: themes?.map((theme) => theme.value) || [],
        enumNames: themes?.map((theme) => theme.label) || [],
      },
      language: {
        type: 'string',
        title: t('language.label'),
        description: t('language.text'),
        enum: langs?.map((lang) => lang.value),
        enumNames: langs?.map((lang) => lang.label),
      },
      time_zone: {
        type: 'string',
        title: t('time_zone.label'),
        description: t('time_zone.text'),
      },
    },
  };

  const [formData, setFormData] = useState<FormDataType>({
    theme: {
      value: setting?.theme || storeInterface.theme,
      isInvalid: false,
      errorMsg: '',
    },
    language: {
      value: setting?.language || storeInterface.language,
      isInvalid: false,
      errorMsg: '',
    },
    time_zone: {
      value: setting?.time_zone || DEFAULT_TIMEZONE,
      isInvalid: false,
      errorMsg: '',
    },
  });

  const uiSchema: UISchema = {
    theme: {
      'ui:widget': 'select',
    },
    language: {
      'ui:widget': 'select',
    },
    time_zone: {
      'ui:widget': 'timezone',
    },
  };
  const getLangs = async () => {
    const res: LangsType[] = await loadLanguageOptions(true);
    setLangs(res);
  };
  // set default theme value
  if (!formData.theme.value && Array.isArray(themes) && themes.length) {
    setFormData({
      ...formData,
      theme: {
        value: themes[0].value,
        isInvalid: false,
        errorMsg: '',
      },
    });
  }

  const checkValidated = (): boolean => {
    let ret = true;
    const { theme, language } = formData;
    const formCheckData = { ...formData };
    if (!theme.value) {
      ret = false;
      formCheckData.theme = {
        value: '',
        isInvalid: true,
        errorMsg: t('theme.msg'),
      };
    }
    if (!language.value) {
      ret = false;
      formCheckData.language = {
        value: '',
        isInvalid: true,
        errorMsg: t('language.msg'),
      };
    }
    setFormData({
      ...formCheckData,
    });
    return ret;
  };
  const onSubmit = (evt: FormEvent) => {
    evt.preventDefault();
    evt.stopPropagation();
    if (checkValidated() === false) {
      return;
    }
    const reqParams: AdminSettingsInterface = {
      theme: formData.theme.value,
      language: formData.language.value,
      time_zone: formData.time_zone.value,
    };

    updateInterfaceSetting(reqParams)
      .then(() => {
        interfaceStore.getState().update(reqParams);
        setupAppLanguage();
        setupAppTimeZone();
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
  // const imgUpload = (file: any) => {
  //   return new Promise((resolve) => {
  //     uploadAvatar(file).then((res) => {
  //       setFormData({
  //         ...formData,
  //         logo: {
  //           value: res,
  //           isInvalid: false,
  //           errorMsg: '',
  //         },
  //       });
  //       resolve(true);
  //     });
  //   });
  // };

  useEffect(() => {
    if (setting) {
      const formMeta = {};
      Object.keys(setting).forEach((k) => {
        formMeta[k] = { ...formData[k], value: setting[k] };
      });
      setFormData({ ...formData, ...formMeta });
    }
  }, [setting]);
  useEffect(() => {
    getLangs();
  }, []);

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
        onSubmit={onSubmit}
        onChange={handleOnChange}
      />
    </>
  );
};

export default Interface;
