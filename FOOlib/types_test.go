package foodlib_test

import (
	foodlib "FOOdBAR/FOOlib"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSortMethodByNumber(t *testing.T) {
	method, err := foodlib.GetSortMethodByNumber(0)
	assert.NoError(t, err)
	assert.Equal(t, foodlib.Inactive, method)

	_, err = foodlib.GetSortMethodByNumber(-1)
	assert.Error(t, err)
}

func TestString2TabType(t *testing.T) {
	assert.Equal(t, foodlib.Invalid, foodlib.String2TabType("InvalidType"))
}

func TestPageData(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()
	tabID := "test-tab1"

	pd1 := foodlib.InitPageData(userID, sessionID, tabID)
	assert.Equal(t, []*foodlib.TabData{}, pd1.TabDatas)

	userID2 := uuid.New()
	sessionID2 := uuid.New()
	tabID2 := "test-tab2"

	pd2 := foodlib.InitPageData(userID2, sessionID2, tabID2)
	assert.Equal(t, []*foodlib.TabData{}, pd2.TabDatas)

	assert.NotEqual(t, pd1, pd2)
	assert.NotEqual(t, pd1.UserID, pd2.UserID)
	assert.NotEqual(t, pd1.SessionID, pd2.SessionID)
	assert.NotEqual(t, pd1.TabID, pd2.TabID)

	td1 := pd1.GetTabDataByType(foodlib.Invalid)

	assert.Nil(t, td1)

	td2 := pd2.GetTabDataByType(foodlib.Recipe)

	assert.Equal(t, true, pd2.IsActive(foodlib.Recipe))
	pd2.SetActive(td2, false)
	assert.Equal(t, false, pd2.IsActive(foodlib.Recipe))

	itemid1 := uuid.New()
	td3 := pd1.GetTabDataByType(foodlib.Pantry)
	_, err := td3.GetTabItem(itemid1)
	assert.Error(t, err)
	ti1 := td3.AddTabItem(&foodlib.TabItem{ItemID: itemid1})
	control, err := td3.GetTabItem(itemid1)
	assert.NotNil(t, ti1)

	updatedparent := pd1.GetTabDataByType(foodlib.Pantry)
	testctl, err := updatedparent.GetTabItem(itemid1)
	assert.Nil(t, err)
	assert.Equal(t, *control, *testctl)
}


func TestTabItemMarshalJSON(t *testing.T) {
	itemID := uuid.New()
	tabItem := &foodlib.TabItem{
		ItemID:   itemID,
		Ttype:    foodlib.Recipe,
		Selected: true,
		Expanded: false,
	}

	data, err := tabItem.MarshalJSON()
	require.NoError(t, err)

	var unmarshalled map[string]interface{}
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, itemID.String(), unmarshalled["item_id"])
	assert.Equal(t, "Recipe", unmarshalled["tab_type"])
	assert.True(t, unmarshalled["selected"].(bool))
	assert.False(t, unmarshalled["expanded"].(bool))
}

func TestTabItemUnmarshalJSON(t *testing.T) {
	itemID := uuid.New()
	data := []byte(`{"item_id": "` + itemID.String() + `", "tab_type": "Recipe", "selected": true, "expanded": false}`)

	var tabItem foodlib.TabItem
	err := tabItem.UnmarshalJSON(data)
	require.NoError(t, err)

	assert.Equal(t, itemID, tabItem.ItemID)
	assert.True(t, tabItem.Selected)
	assert.False(t, tabItem.Expanded)
}

func TestTabDataMarshalJSON(t *testing.T) {
	itemID := uuid.New()
	flippedID := uuid.New()
	tabItem := &foodlib.TabItem{
		ItemID:   itemID,
		Ttype:    foodlib.Recipe,
		Selected: true,
		Expanded: false,
	}
	tabData := &foodlib.TabData{
		Ttype:   foodlib.Recipe,
		Items:   []*foodlib.TabItem{tabItem},
		OrderBy: []int{1, 2, 3},
		Flipped: flippedID,
	}

	data, err := tabData.MarshalJSON()
	require.NoError(t, err)

	var unmarshalled map[string]interface{}
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, "Recipe", unmarshalled["tab_type"])
	assert.Len(t, unmarshalled["items"], 1)
	assert.Equal(t, flippedID.String(), unmarshalled["flipped"])
}

func TestTabDataUnmarshalJSON(t *testing.T) {
	itemID := uuid.New()
	flippedID := uuid.New()
	data := []byte(`{"tab_type": "Recipe", "items": [{"item_id": "` + itemID.String() + `", "tab_type": "Recipe", "selected": true, "expanded": false}], "order_by": [1, 2, 3], "flipped": "` + flippedID.String() + `"}`)

	var tabData foodlib.TabData
	err := tabData.UnmarshalJSON(data)
	require.NoError(t, err)

	require.Len(t, tabData.Items, 1)
	assert.Equal(t, itemID, tabData.Items[0].ItemID)
	assert.True(t, tabData.Items[0].Selected)
	assert.False(t, tabData.Items[0].Expanded)
	assert.Equal(t, []int{1, 2, 3}, tabData.OrderBy)
	assert.Equal(t, flippedID, tabData.Flipped)
}

func TestPageDataMarshalJSON(t *testing.T) {
	sessionID := uuid.New()
	userID := uuid.New()
	tabData := &foodlib.TabData{
		Ttype: foodlib.Recipe,
		Items: nil,
	}
	pageData := &foodlib.PageData{
		SessionID: sessionID,
		UserID:    userID,
		TabID:     "tab-id",
		TabDatas:  []*foodlib.TabData{tabData},
	}

	data, err := pageData.MarshalJSON()
	require.NoError(t, err)

	var unmarshalled map[string]interface{}
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, sessionID.String(), unmarshalled["session_id"])
	assert.Equal(t, userID.String(), unmarshalled["user_id"])
	assert.Equal(t, "tab-id", unmarshalled["tab_id"])
	assert.Len(t, unmarshalled["tab_datas"], 1)
}

func TestPageDataUnmarshalJSON(t *testing.T) {
	sessionID := uuid.New()
	userID := uuid.New()
	data := []byte(`{"session_id": "` + sessionID.String() + `", "user_id": "` + userID.String() + `", "tab_id": "tab-id", "tab_datas": [{"tab_type": "Recipe", "items": [], "order_by": [], "flipped": "00000000-0000-0000-0000-000000000000"}]}`)

	var pageData foodlib.PageData
	err := pageData.UnmarshalJSON(data)
	require.NoError(t, err)

	assert.Equal(t, sessionID, pageData.SessionID)
	assert.Equal(t, userID, pageData.UserID)
	assert.Equal(t, "tab-id", pageData.TabID)
	require.Len(t, pageData.TabDatas, 1)
	assert.Len(t, pageData.TabDatas[0].Items, 0)
	assert.Len(t, pageData.TabDatas[0].OrderBy, 0)
	assert.Equal(t, uuid.Nil, pageData.TabDatas[0].Flipped)
}
